package groot

import (
	"encoding/binary"
	"fmt"
	"os"
)

type unzip_fct func()

type FileReader struct {
	name        string               // path to this file
	f           *os.File             // handle to the raw file
	order       binary.ByteOrder     // file endianness
	nbytes_read uint64               // number of bytes read from this file
	root_dir    Directory            // root directory of this file
	unzipers    map[string]unzip_fct // unziper functions
	title       string               // title of this file

	// -- record --

	version     uint32 // file format version
	beg         int64  // first used byte in file
	end         int64  // last used byte in file
	seek_free   int64  // location on disk of free segments structure
	nbytes_free uint32 // number of bytes for free segments structure
	//nfree int
	nbytes_name uint32 // number of bytes in TNamed at creation time
	seek_info   int64  // location on disk of streamerinfos
	nbytes_info uint32 // number of bytes for streamerinfos?
}

func NewFileReader(name string) (f *FileReader, err error) {
	f = &FileReader{
		name:     name,
		order:    binary.BigEndian,
		unzipers: make(map[string]unzip_fct),
	}

	f.f, err = os.Open(name)
	if err != nil {
		return nil, err
	}

	f.root_dir = Directory{file: f}

	err = f.initialize()
	if err != nil {
		return nil, err
	}
	return f, err
}

func (f *FileReader) breader() breader {
	return breader{f.order}
}

func (f *FileReader) seek(offset int64, pos int) (ret int64, err error) {
	switch pos {
	case os.SEEK_SET:
		return f.f.Seek(offset, os.SEEK_SET)
	case os.SEEK_CUR:
		return f.f.Seek(offset, os.SEEK_CUR)
	case os.SEEK_END:
		return f.f.Seek(offset, os.SEEK_END)

	default:
		err = fmt.Errorf("groot: unknown seek option (%d)", pos)
		panic(err)
	}
	panic("unreachable")
}

func (f *FileReader) initialize() (err error) {
	err = f.read_header()
	if err != nil {
		return err
	}

	cur, err := f.f.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}
	defer f.f.Seek(cur, os.SEEK_SET)

	nbytes := f.nbytes_name + f.root_dir.record_size(f.version)
	printf("nbytes: %v\n", nbytes)

	// read directory info
	_, err = f.f.Seek(f.beg+int64(f.nbytes_name), os.SEEK_SET)
	if err != nil {
		return err
	}
	err = f.root_dir.from_buffer(f.f)
	if err != nil {
		return err
	}

	nk := sz_int32     // Key::fNumberOfBytes
	nk += sz_int16     // Key::fVersion
	nk += 2 * sz_int32 // Key::fObjectSize, Date
	nk += 2 * sz_int16 // Key::fKeyLength, fCycle
	nk += 2 * sz_int32 // Key::fSeekKey, fSeekParentDirectory
	// WARNING: the above is sz_int32 since we are at beginning of file...

	_, err = f.f.Seek(f.beg+int64(nk), os.SEEK_SET)
	if err != nil {
		return err
	}

	br := f.breader()

	var cname string
	cname = br.readTString(f.f)
	if cname != "TFile" {
		return fmt.Errorf("groot: expected [TFile]. got [%v]", cname)
	}
	printf("f-clsname [%v]\n", cname)

	cname = br.readTString(f.f)
	printf("f-cname   [%v]\n", cname)

	f.title = br.readTString(f.f)
	printf("f-title   [%v]\n", f.title)

	if f.root_dir.nbytes_name < 10 || f.root_dir.nbytes_name > 1000 {
		return fmt.Errorf("groot: can't read directory info.")
	}

	// read keys of the top-level directory
	if f.root_dir.seek_keys <= f.beg {
		return fmt.Errorf("groot: file [%s] is probably not closed", f.name)
	}
	nkeys, err := f.root_dir.read_keys()
	if err != nil {
		return err
	}
	printf("f-dir-nkeys: %v\n", nkeys)
	return nil
}

func (f *FileReader) read_header() (err error) {
	cur, err := f.f.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}
	defer f.f.Seek(cur, os.SEEK_SET)

	_, err = f.f.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	br := f.breader()
	{
		hdr := make([]byte, 4)
		err = binary.Read(f.f, f.order, &hdr)
		if err != nil {
			return err
		}
		printf("hdr: %v\n", string(hdr))
		if string(hdr) != "root" {
			return fmt.Errorf(
				"groot: file [%s] is not a ROOT file (%v)",
				f.name, string(hdr))
		}
	}
	f.version = br.ntou4(f.f)
	f.beg = int64(br.ntou4(f.f))
	printf("beg: %v\n", f.beg)
	if f.version >= 1000000 {
		f.end = int64(br.ntou8(f.f))
		f.seek_free = int64(br.ntou8(f.f))
	} else {
		f.end = int64(br.ntou4(f.f))
		f.seek_free = int64(br.ntou4(f.f))
	}
	printf("end: %v\n", f.end)
	printf("seek-free: %v\n", f.seek_free)
	f.nbytes_free = br.ntou4(f.f)
	/*nfree*/ br.ntoi4(f.f)
	f.nbytes_name = br.ntou4(f.f)
	printf("nbytes-free: %v\n", f.nbytes_free)
	printf("nbytes-name: %v\n", f.nbytes_name)
	/*units*/ br.ntobyte(f.f)
	/*compress*/ br.ntou4(f.f)
	if f.version >= 1000000 {
		f.seek_info = int64(br.ntou8(f.f))
	} else {
		f.seek_info = int64(br.ntou4(f.f))
	}
	f.nbytes_info = br.ntou4(f.f)
	printf("seek-info: %v\n", f.seek_info)
	printf("nbytes-info: %v\n", f.nbytes_info)

	// read streamer infos
	return f.read_streamer_infos()
}

func (f *FileReader) read_streamer_infos() (err error) {
	buf := make([]byte, int(f.nbytes_info))
	_, err = f.f.ReadAt(buf, f.seek_info)
	if err != nil {
		return err
	}
	printf("buf: %v\n", len(buf))
	return err
}

func (f *FileReader) Name() string {
	return f.name
}

func (f *FileReader) Version() uint32 {
	return f.version
}

func (f *FileReader) Dir() *Directory {
	return &f.root_dir
}

func (f *FileReader) ByteOrder() binary.ByteOrder {
	return f.order
}

// EOF
