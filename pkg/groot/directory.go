package groot

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

// Directory is a directory inside a ROOT file
type Directory struct {
	file        *File
	keys        []Key
	ctime       time.Time // time of directory's creation
	mtime       time.Time // time of directory's last modification
	nbytes_keys uint32    // number of bytes for the keys
	nbytes_name uint32    // number of bytes in TNamed at creation time
	seek_dir    int64     // location of directory on file
	seek_parent int64     // location of parent directory on file
	seek_keys   int64     // location of keys record on file
}

func NewDirectory(f *File, buf []byte) (d *Directory, err error) {
	d = &Directory{file:f}
	b, err := NewBuffer(buf, f.order, 0)
	err = d.from_buffer(b)
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

func (d *Directory) from_buffer(b *Buffer) (err error) {
	version := b.ntou2()
	d.ctime = datime2time(b.ntou4())
	d.mtime = datime2time(b.ntou4())
	d.nbytes_keys = b.ntou4()
	d.nbytes_name = b.ntou4()
	if version > 1000 {
		d.seek_dir = b.ntoi8()
		d.seek_parent = b.ntoi8()
		d.seek_keys = b.ntoi8()
	} else {
		d.seek_dir = int64(b.ntoi4())
		d.seek_parent = int64(b.ntoi4())
		d.seek_keys = int64(b.ntoi4())
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

func (d *Directory) ROOTDecode(b *Buffer) (err error) {
	version := b.ntou2()
	d.ctime = datime2time(b.ntou4())
	d.mtime = datime2time(b.ntou4())
	d.nbytes_keys = b.ntou4()
	d.nbytes_name = b.ntou4()
	if version > 1000 {
		d.seek_dir = b.ntoi8()
		d.seek_parent = b.ntoi8()
		d.seek_keys = b.ntoi8()
	} else {
		d.seek_dir = int64(b.ntoi4())
		d.seek_parent = int64(b.ntoi4())
		d.seek_keys = int64(b.ntoi4())
	}
	printf("dir-version: %v\n", version)
	printf("dir-ctime: %v\n", d.ctime.String())
	printf("dir-mtime: %v\n", d.mtime.String())
	printf("dir-nbytes-keys: %v\n", d.nbytes_keys)
	printf("dir-nbytes-name: %v\n", d.nbytes_name)
	printf("dir-seek-dir: %v\n", d.seek_dir)
	printf("dir-seek-parent: %v\n", d.seek_parent)
	printf("dir-seek-keys: %v\n", d.seek_keys)

	if err != nil {
		return err
	}
	_, err = d.read_keys()
	return err
}

func (d *Directory) ROOTEncode(buffer *Buffer) error {
	panic("groot.Directory.ROOTEncode: sorry, not implemented")
}

func (d *Directory) SetFile(f *File) error {
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

	b,err := NewBuffer(buf, d.file.order, 0)
	if err != nil {
		return -1, err
	}

	err = hdr.init_from_buffer(b)
	if err != nil {
		return -1, err
	}

	nkeys = int(b.ntoi4())
	printf("dir-nkeys: %v\n", nkeys)

	d.keys = make([]Key, nkeys)
	for i := 0; i < nkeys; i++ {
		printf("--key-- %v\n", i)
		key, err := NewKey(d.file, 0, 0)
		if err != nil {
			return -1, err
		}
		err = key.init_from_buffer(b)
		if err != nil {
			return -1, err
		}
		d.keys[i] = *key
	}
	return nkeys, nil
}

func init() {

	new_dir := func() reflect.Value {
		o := &Directory{file: nil, keys: make([]Key, 0)}
		return reflect.ValueOf(o)
	}

	Factory.db["TDirectory"] = new_dir
	Factory.db["*groot.Directory"] = new_dir
}


// EOF
