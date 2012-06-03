package groot

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	sz_int16  = 2
	sz_int32  = 4
	sz_int64  = 8
	sz_uint16 = 2
	sz_uint32 = 4
	sz_uint64 = 8

	g_START_BIG_FILE = 2000000000
)

// Class represents a ROOT class.
// Class instances are created by a ClassFactory.
type Class interface {
	// GetCheckSum gets the check sum for this ROOT class
	GetCheckSum() int

	// GetMembers returns the list of members for this ROOT class
	GetMembers() []Member

	// GetVersion returns the version number for this ROOT class
	GetVersion() int

	// GetClassName returns the ROOT class name for this ROOT class
	GetClassName() string

	// GetSuperClasses returns the list of super-classes for this ROOT class
	GetSuperClasses() []Class
}

// Member represents a single member of a ROOT class
type Member interface {
	// GetArrayDim returns the dimension of the array (if any)
	GetArrayDim() int

	// GetComment returns the comment associated with this member
	GetComment() string

	// GetName returns the name of this member
	GetName() string

	// GetType returns the class of this member
	GetType() Class

	// GetValue returns the value of this member
	//GetValue(o Object) reflect.Value
}

// Object represents a ROOT object
type Object interface {
	// GetClass returns the ROOT class of this object
	GetClass() Class
}

// ClassFactory creates ROOT classes
type ClassFactory interface {
	Create(name string) Class
}

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
	println("dir-version:", version)
	println("dir-ctime:", d.ctime.String())
	println("dir-mtime:", d.mtime.String())
	println("dir-nbytes-keys:", d.nbytes_keys)
	println("dir-nbytes-name:", d.nbytes_name)
	println("dir-seek-dir:", d.seek_dir)
	println("dir-seek-parent:", d.seek_parent)
	println("dir-seek-keys:", d.seek_keys)
	return err
}

// read_keys reads the keys from a Directory
//
// Every Directory has a list of keys (fKeys).
// This list has been written on the file via ROOT::TDirectory::writeKeys
// as a single data record.
func (d *Directory) read_keys() (nkeys int, err error) {

	println("--read_keys--")
	hdr, err := NewKey(d.file, d.seek_keys, d.nbytes_keys)
	if err != nil {
		return -1, err
	}
	if hdr == nil {
		return -1, fmt.Errorf("groot: invalid header key")
	}

	cur, err := d.file.f.Seek(0, os.SEEK_CUR)
	if err != nil {
		return -1, err
	}
	defer d.file.f.Seek(cur, os.SEEK_SET)

	buf := make([]byte, int(d.nbytes_keys))

	_, err = d.file.f.ReadAt(buf, d.seek_keys)
	if err != nil {
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
	println("dir-nkeys:", nkeys)

	d.keys = make([]Key, nkeys)
	for i := 0; i < nkeys; i++ {
		println("--key--", i)
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
	println("key-nbytes:", k.nbytes)

	println("key-version:", k.version)

	k.version = uint32(br.ntou2(f))
	println("key-version:", k.version)
	k.objsz = uint32(br.ntoi4(f))
	println("key-objsz:", k.objsz)

	k.date = datime2time(br.ntou4(f))
	println("key-cdate:", k.date.String())

	k.keysz = br.ntou2(f)
	println("key-keysz:", k.keysz)

	k.cycle = br.ntou2(f)
	println("key-cycle:", k.cycle)

	if k.version > 1000 {
		k.seek_key = br.ntoi8(f)
		k.seek_parent_dir = br.ntoi8(f)
	} else {
		k.seek_key = int64(br.ntoi4(f))
		k.seek_parent_dir = int64(br.ntoi4(f))
	}
	println("key-seek-key:", k.seek_key)
	println("key-seek-pdir:", k.seek_parent_dir)

	k.class = br.readTString(f)
	println("key-class [" + k.class + "]")

	k.name = br.readTString(f)
	println("key-name  [" + k.name + "]")

	k.title = br.readTString(f)
	println("key-title [" + k.title + "]")

	return err
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
