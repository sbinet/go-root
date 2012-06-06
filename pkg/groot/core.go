package groot

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

var g_verbose = false
func printf(format string, args ...interface{}) {
	if g_verbose {
		fmt.Printf(format, args...)
	}
}

const (
	sz_int16  = 2
	sz_int32  = 4
	sz_int64  = 8
	sz_uint16 = 2
	sz_uint32 = 4
	sz_uint64 = 8

	g_START_BIG_FILE = 2000000000
)

// Directory is a directory inside a ROOT file
type Directory struct {
	file        *FileReader
	keys        []Key
	ctime       time.Time // time of directory's creation
	mtime       time.Time // time of directory's last modification
	nbytes_keys uint32    // number of bytes for the keys
	nbytes_name uint32    // number of bytes in TNamed at creation time
	seek_dir    int64     // location of directory on file
	seek_parent int64     // location of parent directory on file
	seek_keys   int64     // location of keys record on file
}

func NewDirectory(f *FileReader, buf []byte) (d *Directory, err error) {
	d = &Directory{file:f}
	err = d.from_buffer(bytes.NewBuffer(buf))
	if err != nil {
		return
	}
	_, err = d.read_keys()
	return
}

func (d *Directory) Keys() []Key {
	return d.keys
}

func (d *Directory) record_size(version uint32) uint32 {
	var nbytes uint32 = sz_uint16
	nbytes += sz_uint32 // ctime
	nbytes += sz_uint32 // mtime
	nbytes += sz_uint32 // nbytes_keys
	nbytes += sz_uint32 // nbytes_name
	if version >= 40000 {
		nbytes += sz_int64 // seek_dir
		nbytes += sz_int64 // seek_parent
		nbytes += sz_int64 // seek_keys
	} else {
		nbytes += sz_int32 // seek_dir
		nbytes += sz_int32 // seek_parent
		nbytes += sz_int32 // seek_keys
	}
	return nbytes
}

func (d *Directory) from_buffer(buf io.Reader) (err error) {
	br := d.file.breader()
	version := br.ntou2(buf)
	d.ctime = datime2time(br.ntou4(buf))
	d.mtime = datime2time(br.ntou4(buf))
	d.nbytes_keys = br.ntou4(buf)
	d.nbytes_name = br.ntou4(buf)
	if version > 1000 {
		d.seek_dir = br.ntoi8(buf)
		d.seek_parent = br.ntoi8(buf)
		d.seek_keys = br.ntoi8(buf)
	} else {
		d.seek_dir = int64(br.ntoi4(buf))
		d.seek_parent = int64(br.ntoi4(buf))
		d.seek_keys = int64(br.ntoi4(buf))
	}
	printf("dir-version: %v\n", version)
	printf("dir-ctime: %v\n", d.ctime.String())
	printf("dir-mtime: %v\n", d.mtime.String())
	printf("dir-nbytes-keys: %v\n", d.nbytes_keys)
	printf("dir-nbytes-name: %v\n", d.nbytes_name)
	printf("dir-seek-dir: %v\n", d.seek_dir)
	printf("dir-seek-parent: %v\n", d.seek_parent)
	printf("dir-seek-keys: %v\n", d.seek_keys)
	return err
}

func (d *Directory) ROOTDecode(buf []byte) (err error) {
	iobuf := bytes.NewBuffer(buf)
	err = d.from_buffer(iobuf)
	if err != nil {
		return err
	}
	_, err = d.read_keys()
	return err
}

func (d *Directory) ROOTEncode(buf []byte) error {
	panic("groot.Directory.ROOTEncode: sorry, not implemented")
}

func (d *Directory) SetFile(f *FileReader) error {
	d.file = f
	return nil
}

// read_keys reads the keys from a Directory
//
// Every Directory has a list of keys (fKeys).
// This list has been written on the file via ROOT::TDirectory::writeKeys
// as a single data record.
func (d *Directory) read_keys() (nkeys int, err error) {

	printf("--read_keys-- %v %v\n", d.seek_keys, d.nbytes_keys)
	hdr, err := NewKey(d.file, d.seek_keys, d.nbytes_keys)
	if err != nil {
		printf("groot.Directory.read_keys: %v\n",err.Error())
		return -1, err
	}
	if hdr == nil {
		return -1, fmt.Errorf("groot: invalid header key")
	}

	cur, err := d.file.f.Seek(0, os.SEEK_CUR)
	if err != nil {
		printf("groot.Directory.read_keys: %v\n",err.Error())
		return -1, err
	}
	defer d.file.f.Seek(cur, os.SEEK_SET)

	buf := make([]byte, int(d.nbytes_keys))

	printf("--- %v %v\n", len(buf), d.seek_keys)
	_, err = d.file.f.ReadAt(buf, d.seek_keys)
	if err != nil {
		printf("seek_keys: %v\n",d.seek_keys)
		printf("len(buf): %v\n",len(buf))
		printf("groot.Directory.read_keys-ReadAt: %v\n",err.Error())
		return -1, err
	}

	f := bytes.NewReader(buf)
	if f == nil {
		return -1, fmt.Errorf("groot: could not create a bytes.Reader")
	}

	err = hdr.init_from_buffer(f)
	if err != nil {
		return -1, err
	}

	br := d.file.breader()
	nkeys = int(br.ntoi4(f))
	printf("dir-nkeys: %v\n", nkeys)

	d.keys = make([]Key, nkeys)
	for i := 0; i < nkeys; i++ {
		printf("--key-- %v\n", i)
		key, err := NewKey(d.file, 0, 0)
		if err != nil {
			return -1, err
		}
		err = key.init_from_buffer(f)
		if err != nil {
			return -1, err
		}
		d.keys[i] = *key
	}
	return nkeys, nil
}

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
	file *FileReader
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

func NewKey(f *FileReader, pos int64, nbytes uint32) (k *Key, err error) {
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
