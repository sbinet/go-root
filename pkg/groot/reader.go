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
	r             io.ReadSeeker
	hdr           fileHeader
	url           string

	dirs []Directory
	diridx int
	
	keys []Key
	keyidx int

	buffers []string
	bufidx int

	//hists []Histogram
	//histidx int

	seekinfo uint64
	nbytesinfo uint64

	streamers []string
	streamerinfo string

	offset        uint64
	archiveoffset uint64
	name          string
	title         string
	dir           Directory

	datimeC       time.Time
	datimeM       time.Time
	fileClass     Class
	factory       ClassFactory

	streamerInfo  Key
}

type fileHeader struct {
	version    uint32
	beg        uint32
	end        uint64
	units      byte
	seekinfo   uint64
	nbytesinfo uint64
}

type Key struct {
	offset uint64
	nbytes uint64
	version uint16
	objlen uint32
	datime uint32 //fixme: use time.Time ?
	keylen uint16
	cycle  uint16
	seekkey  uint64
	seekpdir uint64
	class  string
	name string
	title string
}

func (k *Key) dataoffset() uint64 {
	return k.seekkey + uint64(k.keylen)
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

func (b breader) readKey(r io.ReadSeeker) Key {
	key := Key{}
	offset, err := r.Seek(0, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	key.offset = uint64(offset)

	nbytes := b.ntoi4(r)
	if nbytes < 0 {
		key.nbytes = -uint64(nbytes)
	} else {
		key.nbytes = +uint64(nbytes)
	}
	key.version = b.ntou2(r)
	key.objlen = b.ntou4(r)
	println("nbytes:", key.nbytes)
	println("version:", key.version)
	println("objlen:", key.objlen)

	key.datime = b.ntou4(r)
	println("datime:", key.datime)

	{
		var year uint32 = (key.datime >> 26) + 1995
		var month uint32 = (key.datime << 6) >> 28
		var day uint32 = (key.datime << 10) >> 27
		var hour uint32 = (key.datime << 15) >> 27
		var min uint32 = (key.datime << 20) >> 26
		var sec uint32 = (key.datime << 26) >> 26
		println(year, month, day, hour, min, sec)
	}
	key.keylen = b.ntou2(r)
	key.cycle = b.ntou2(r)
	println("keylen:", key.keylen)
	println("cycle:", key.cycle)

	largeKey := int64(key.nbytes)+1 > 2*1024*1024*1024 /*2G*/
	if largeKey {
		key.seekkey = b.ntou8(r)
		key.seekpdir = b.ntou8(r)
	} else {
		key.seekkey = uint64(b.ntou4(r))
		key.seekpdir = uint64(b.ntou4(r))
	}
	key.class = b.readTString(r)
	key.name = b.readTString(r)
	key.title = b.readTString(r)
	println("classname:", key.class)
	println("fname:", key.name)
	println("title:", key.title)

	println("dataoffset:", key.dataoffset())
	pos, err := r.Seek(0, os.SEEK_CUR)
	println("key-offset:", key.offset)
	println("current-offset:", pos)
	return key
}

func (b breader) readKeys(r io.ReadSeeker) []Key {
	keys := make([]Key, 0)
	pos, _ := r.Seek(0, os.SEEK_CUR)
	defer r.Seek(pos, os.SEEK_SET)

	//cbk1 := func() {
	//}

	return keys
}

func NewFileReader(name string) (f *FileReader, err error) {
	f = &FileReader{}
	f.url = name

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
	f.hdr = br.readHeader(f.r)
	println("version:", f.hdr.version)

	if f.hdr.version < 30006 {
		return nil, fmt.Errorf("groot: cannot read ROOT files created with version <= 3.00.06")
	}

	println("version:", f.hdr.version)
	println("beg=", f.hdr.beg)

	println("seekinfo:", f.hdr.seekinfo)
	_, err = f.r.Seek(int64(f.hdr.beg), os.SEEK_SET)
	if err != nil {
		return nil, err
	}

	key := br.readKey(f.r)
	pos, err := f.r.Seek(0, os.SEEK_CUR)
	println("current-offset:", pos)
	println("key:",key.name)
	f.name = key.name
	f.title = key.title
	return f, err
}

func (f *FileReader) GetName() string {
	return f.name
}

func (f *FileReader) GetTitle() string {
	return f.title
}

func (f *FileReader) GetVersion() uint32 {
	return f.hdr.version
}

// EOF
