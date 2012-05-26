package groot

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
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

//
type FileReader struct {
	datimeC      time.Time
	datimeM      time.Time
	fileClass    Class
	factory      ClassFactory
	r            io.ReadSeeker
	name         string
	title        string
	dir          Directory
	streamerInfo Key
	nbytesKeys   int32
	nbytesName   int32
	seekDir      int64
	seekKeys     int64
	seekParent   int64
	version      uint32
	seekInfo     uint64
}

type fileHeader struct {
	version uint32
	beg uint32
	end uint64
	units byte
	seekinfo uint64
	nbytesinfo uint64
}

func (b breader) readHeader(r io.Reader) fileHeader {
	hdr := fileHeader{}
	hdr.version = b.ntou4(r)
	hdr.beg = b.ntou4(r)
	large_file := hdr.version >= 1000000
	if large_file {
		hdr.end = b.ntou8(r)
		hdr.units = b.ntobyte(r)
		hdr.seekinfo = b.ntou8(r)
		hdr.nbytesinfo = b.ntou8(r)
	} else {
		hdr.end = uint64(b.ntou4(r))
		hdr.units = b.ntobyte(r)
		hdr.seekinfo = uint64(b.ntou4(r))
		hdr.nbytesinfo = uint64(b.ntou4(r))
	}
	return hdr
}

func NewFileReader(name string) (f *FileReader, err error) {
	f = &FileReader{}
	f.r, err = os.Open(name)
	if err != nil {
		return nil, err
	}

	order := binary.BigEndian
	br := breader{order}

	{
		hdr := make([]byte, 4)
		err = binary.Read(f.r, order, &hdr)
		if err != nil {
			return nil, err
		}
		println("hdr:", string(hdr))
		if string(hdr) != "root" {
			return nil, fmt.Errorf("groot: file [%s] is not a ROOT file (%v)",
				name, string(hdr))
		}
	}
	hdr := br.readHeader(f.r)
	f.version = hdr.version
	f.seekInfo = hdr.seekinfo
	println("version:", f.version)

	if f.version < 30006 {
		return nil, fmt.Errorf("groot: cannot read ROOT files created with version <= 3.00.06")
	}

	println("version:", f.version)
	println("beg=", hdr.beg)


	println("seekinfo:", f.seekInfo)
	_, err = f.r.Seek(int64(hdr.beg), 0)
	if err != nil {
		return nil, err
	}

	nbytes := br.ntoi4(f.r)
	if nbytes < 0 {
		nbytes = -nbytes
	}
	version := br.ntou2(f.r)
	objlen := br.ntou4(f.r)
	println("nbytes:", nbytes)
	println("version:", version)
	println("objlen:", objlen)

	datime := br.ntou4(f.r)
	println("datime:", datime)

	{
		var year uint32 = (datime >> 26) + 1995
		var month uint32 = (datime << 6) >> 28
		var day uint32 = (datime << 10) >> 27
		var hour uint32 = (datime << 15) >> 27
		var min uint32 = (datime << 20) >> 26
		var sec uint32 = (datime << 26) >> 26
		println(year, month, day, hour, min, sec)
	}
	keylen := br.ntou2(f.r)
	cycle := br.ntou2(f.r)
	println("keylen:",keylen)
	println("cycle:",cycle)

	largeKey := int64(nbytes)+1 > 2 * 1024 * 1024 * 1024 /*2G*/
	var seekkey uint64
	if largeKey {
		seekkey = br.ntou8(f.r)
		_ = br.ntou8(f.r) // skip seekPdir
	} else {
		seekkey = uint64(br.ntou4(f.r))
		_ = br.ntou4(f.r) // skip seekPdir
	}
	classname := br.readTString(f.r)
	f.name = br.readTString(f.r)
	f.title = br.readTString(f.r)
	println("classname:",classname)
	println("fname:",f.name)
	println("title:",f.title)

	dataoffset := seekkey + uint64(keylen)
	println("dataoffset:",dataoffset)
	return f, err
}

func (f *FileReader) GetName() string {
	return f.name
}

func (f *FileReader) GetTitle() string {
	return f.title
}

func (f *FileReader) GetVersion() uint32 {
	return f.version
}

// EOF
