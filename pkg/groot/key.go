package groot

import (
	"io"
	"time"
)

// Key is a key (a label) in a ROOT file
//
//  The Key class includes functions to book space on a file,
//   to create I/O buffers, to fill these buffers
//   to compress/uncompress data buffers.
//
//  Before saving (making persistent) an object on a file, a key must
//  be created. The key structure contains all the information to
//  uniquely identify a persistent object on a file.
//     fNbytes    = number of bytes for the compressed object+key
//     version of the Key class
//     fObjlen    = Length of uncompressed object
//     fDatime    = Date/Time when the object was written
//     fKeylen    = number of bytes for the key structure
//     fCycle     = cycle number of the object
//     fSeekKey   = Address of the object on file (points to fNbytes)
//                  This is a redundant information used to cross-check
//                  the data base integrity.
//     fSeekPdir  = Pointer to the directory supporting this object
//     fClassName = Object class name
//     fName      = name of the object
//     fTitle     = title of the object
//
//  The Key class is used by ROOT to:
//    - to write an object in the Current Directory
//    - to write a new ntuple buffer
type Key struct {
	file *File
	//bufsz  uint32
	//buffer []byte

	// record -- stored in file

	nbytes          uint32    // number of bytes for the object on file
	version         uint32    // key version identifier
	objsz           uint32    // length of uncompressed object in bytes
	date            time.Time // time of insertion in file
	keysz           uint16    // number of bytes for the key itself
	cycle           uint16    // cycle number
	seek_key        int64     // location of object on file
	seek_parent_dir int64     // location of parent directory on file
	class           string    // object class name
	name            string    // name of the object
	title           string    // title of the object
}

func NewKey(f *File, pos int64, nbytes uint32) (k *Key, err error) {
	k = &Key{
		file:     f,
		seek_key: pos,
		nbytes:   nbytes,
		version:  2,
	}
	if pos > g_START_BIG_FILE {
		k.version += 1000
	}
	return k, err
}

func (k *Key) init_from_buffer(f io.Reader) (err error) {

	br := k.file.breader()

	// read the key structure from the buffer
	k.nbytes = br.ntou4(f)
	printf("key-nbytes: %v\n", k.nbytes)

	printf("key-version: %v\n", k.version)

	k.version = uint32(br.ntou2(f))
	printf("key-version: %v\n", k.version)
	k.objsz = uint32(br.ntoi4(f))
	printf("key-objsz: %v\n", k.objsz)

	k.date = datime2time(br.ntou4(f))
	printf("key-cdate: %v\n", k.date.String())

	k.keysz = br.ntou2(f)
	printf("key-keysz: %v\n", k.keysz)

	k.cycle = br.ntou2(f)
	printf("key-cycle: %v\n", k.cycle)

	if k.version > 1000 {
		k.seek_key = br.ntoi8(f)
		k.seek_parent_dir = br.ntoi8(f)
	} else {
		k.seek_key = int64(br.ntoi4(f))
		k.seek_parent_dir = int64(br.ntoi4(f))
	}
	printf("key-seek-key: %v\n", k.seek_key)
	printf("key-seek-pdir: %v\n", k.seek_parent_dir)

	k.class = br.readTString(f)
	printf("key-class [%v]\n", k.class)

	k.name = br.readTString(f)
	printf("key-name  [%v]\n", k.name)

	k.title = br.readTString(f)
	printf("key-title [%v]\n", k.title)

	return err
}

// Buffer returns the buffer of bytes corresponding to the Key's value
func (k *Key) Buffer() (buf []byte, err error) {
	buf = make([]byte, 0)

	if k.keysz == 0 {
		printf("groot.Key.Buffer: key size is zero\n")
		return
	}

	if k.nbytes == 0 {
		printf("groot.Key.Buffer: nbytes is zero\n")
		return
	}

	printf("--Key.Buffer--\n")
	printf("nbytes: %v\n",k.nbytes)
	printf("keysz: %v\n",k.keysz)
	printf("objsz: %v\n",k.objsz)
	printf("seek-key: %v\n",k.seek_key)
	printf("compressed: %v\n", (k.objsz > (k.nbytes - uint32(k.keysz))))

	if k.objsz <= (k.nbytes - uint32(k.keysz)) {
		bufsz := int(k.nbytes - uint32(k.keysz))
		if bufsz < int(k.nbytes) {
			bufsz = int(k.nbytes)
		}
		buf = make([]byte, bufsz)
		printf("*** %v %v\n",len(buf), k.seek_key)
		_, err = k.file.f.ReadAt(buf, k.seek_key)
		if err != nil {
			return []byte{}, err
		}
		// extract the pure object-buffer
		buf = buf[k.keysz:]
	} else {
		// have to decompress
		// size of compressed buffer
		compsz := int(k.objsz + uint32(k.keysz))
		compbuf := make([]byte, compsz)

		_, err = k.file.f.ReadAt(compbuf, k.seek_key)
		if err != nil {
			return []byte{}, err
		}
		buf, err = unzip_root_buffer(compbuf[k.keysz:])
		if err != nil {
			return []byte{}, err
		}
	}
	return
}

func (k *Key) Value() (v interface{}) {
	factory := Factory.Get(k.Class())
	if factory == nil {
		return v
	}

	vv := factory()
	if vv,ok := vv.Interface().(FileSetter); ok {
		err := vv.SetFile(k.file)
		if err != nil {
			return v
		}
	}
	if vv,ok := vv.Interface().(ROOTStreamer); ok {
		buf, err := k.Buffer()
		if err != nil {
			return v
		}
		err = vv.ROOTDecode(buf)
		if err != nil {
			return v
		}
	}
	v = vv.Interface()
	return v
}

func (k *Key) Size() uint32 {
	return k.objsz
}

func (k *Key) Class() string {
	return k.class
}

func (k *Key) Name() string {
	return k.name
}

func (k *Key) Title() string {
	return k.title
}

// EOF
