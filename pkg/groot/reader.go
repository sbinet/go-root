package groot

import (
	"encoding/binary"
	"fmt"
	"os"
)

// constants for the streamers
const (
	kBase       = 0
	kChar       = 1
	kShort      = 2
	kInt        = 3
	kLong       = 4
	kFloat      = 5
	kCounter    = 6
	kCharStar   = 7
	kDouble     = 8
	kDouble32   = 9
	kLegacyChar = 10
	kUChar      = 11
	kUShort     = 12
	kUInt       = 13
	kULong      = 14
	kBits       = 15
	kLong64     = 16
	kULong64    = 17
	kBool       = 18
	kFloat16    = 19
	kOffsetL    = 20
	kOffsetP    = 40
	kObject     = 61
	kAny        = 62
	kObjectp    = 63
	kObjectP    = 64
	kTString    = 65
	kTObject    = 66
	kTNamed     = 67
	kAnyp       = 68
	kAnyP       = 69
	kAnyPnoVT   = 70
	kSTLp       = 71

	kSkip  = 100
	kSkipL = 120
	kSkipP = 140

	kConv  = 200
	kConvL = 220
	kConvP = 240

	kSTL       = 300
	kSTLstring = 365

	kStreamer   = 500
	kStreamLoop = 501
)

const kByteCountMask = 0x40000000

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

	f.root_dir = Directory{file:f}

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
	println("nbytes:", nbytes)
	
	// read directory info
	_, err = f.f.Seek(f.beg + int64(f.nbytes_name), os.SEEK_SET)
	if err != nil {
		return err
	}
	err = f.root_dir.from_buffer(f.f)
	if err != nil {
		return err
	}

	nk := sz_int32 // Key::fNumberOfBytes
	nk += sz_int16 // Key::fVersion
	nk += 2*sz_int32 // Key::fObjectSize, Date
	nk += 2*sz_int16 // Key::fKeyLength, fCycle
	nk += 2*sz_int32 // Key::fSeekKey, fSeekParentDirectory
	// WARNING: the above is sz_int32 since we are at beginning of file...

	
	_, err = f.f.Seek(f.beg + int64(nk), os.SEEK_SET)
	if err != nil {
		return err
	}

	br := f.breader()

	var cname string
	cname = br.readTString(f.f)
	if cname != "TFile" {
		return fmt.Errorf("groot: expected [TFile]. got [%v]", cname)
	}
	println("f-clsname ["+cname+"]")

	cname = br.readTString(f.f)
	println("f-cname   ["+cname+"]")

	f.title = br.readTString(f.f)
	println("f-title   ["+f.title+"]")

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
	println("f-dir-nkeys:", nkeys)
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
		println("hdr:", string(hdr))
		if string(hdr) != "root" {
			return fmt.Errorf(
				"groot: file [%s] is not a ROOT file (%v)",
				f.name, string(hdr))
		}
	}
	f.version = br.ntou4(f.f)
	f.beg = int64(br.ntou4(f.f))
	println("beg:", f.beg)
	if f.version >= 1000000 {
		f.end = int64(br.ntou8(f.f))
		f.seek_free = int64(br.ntou8(f.f))
	} else {
		f.end = int64(br.ntou4(f.f))
		f.seek_free = int64(br.ntou4(f.f))
	}
	println("end:", f.end)
	println("seek-free:", f.seek_free)
	f.nbytes_free = br.ntou4(f.f)
	/*nfree*/ br.ntoi4(f.f)
	f.nbytes_name = br.ntou4(f.f)
	println("nbytes-free:", f.nbytes_free)
	println("nbytes-name:", f.nbytes_name)
	return nil
}

func (f *FileReader) Name() string {
	return f.name
}

func (f *FileReader) Version() uint32 {
	return f.version
}

// EOF
