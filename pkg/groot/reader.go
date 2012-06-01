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
	r   io.ReadSeeker
	order binary.ByteOrder
	hdr fileHeader
	url string

	dirs    []Directory
	keys    []Key
	buffers []string

	//hists []Histogram
	//histidx int

	seekinfo   uint64
	nbytesinfo uint64

	streamers    []string
	streamerinfo string

	offset        uint64
	archiveoffset uint64
	name          string
	title         string
	dir           Directory

	datimeC   time.Time
	datimeM   time.Time
	fileClass Class
	factory   ClassFactory

	streamerInfo Key
}

type fileHeader struct {
	version    uint32
	beg        uint32
	end        uint64
	seekfree   uint64
	nbytesfree uint32
	nfree      uint32
	nbytesname uint32
	units      byte
	compress   uint32
	seekinfo   uint64
	nbytesinfo uint64
	uuid       [18]byte
}

type Key struct {
	offset   uint64
	nbytes   uint64
	version  uint16
	objlen   uint32
	datime   time.Time
	keylen   uint16
	cycle    uint16
	seekkey  uint64
	seekpdir uint64
	class    string
	name     string
	title    string
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
		hdr.seekfree = b.ntou8(r)
		hdr.nbytesfree = b.ntou4(r)
		hdr.nfree = b.ntou4(r)
		hdr.nbytesname = b.ntou4(r)
		hdr.units = b.ntobyte(r)
		hdr.compress = b.ntou4(r)
		hdr.seekinfo = b.ntou8(r)
		hdr.nbytesinfo = b.ntou8(r)
		err := binary.Read(r, b.order, &hdr.uuid)
		if err != nil {
			panic(err)
		}
	} else {
		hdr.end = uint64(b.ntou4(r))
		hdr.seekfree = uint64(b.ntou4(r))
		hdr.nbytesfree = b.ntou4(r)
		hdr.nfree = uint32(b.ntoi4(r))
		hdr.nbytesname = b.ntou4(r)
		hdr.units = b.ntobyte(r)
		hdr.compress = b.ntou4(r)
		hdr.seekinfo = uint64(b.ntou4(r))
		hdr.nbytesinfo = uint64(b.ntou4(r))
		err := binary.Read(r, b.order, &hdr.uuid)
		if err != nil {
			panic(err)
		}
	}
	return hdr
}

func datime2time(d uint32) time.Time {

	// ROOT's TDatime begins in January 1995...
	var year uint32 = (d >> 26) + 1995
	var month uint32 = (d << 6) >> 28
	var day uint32 = (d << 10) >> 27
	var hour uint32 = (d << 15) >> 27
	var min uint32 = (d << 20) >> 26
	var sec uint32 = (d << 26) >> 26
	nsec := 0
	return time.Date(int(year), time.Month(month), int(day),
		int(hour), int(min), int(sec), nsec, time.UTC)
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

	key.datime = datime2time(b.ntou4(r))
	println("datime:", key.datime.String())
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
	println("classname:", key.class)
	key.name = b.readTString(r)
	println("fname:", key.name)
	key.title = b.readTString(r)
	println("title:", key.title)

	println("dataoffset:", key.dataoffset())
	pos, err := r.Seek(0, os.SEEK_CUR)
	println("key-offset:", key.offset)
	println("current-offset:", pos)
	return key
}

func (b breader) readKeys(f *FileReader) []Key {
	keys := make([]Key, 0)
	pos, _ := f.r.Seek(0, os.SEEK_CUR)
	defer f.seek(pos, os.SEEK_SET)

	println("init-pos:",pos)

	cbk1 := func(init int64, f *FileReader) {
		println("--cbk1--")
		println("pos:",init)
		f.seek(init, os.SEEK_SET)
		br := f.breader()
		hdr := br.readHeader(f.r)
		println("header:",hdr.beg, hdr.end)
		cbk2 := func(init int64, f *FileReader) {
			println("--cbk2--")
			br := f.breader()
			f.seek(init, os.SEEK_SET) // skip the "root" identifier
			ver := br.ntou4(f.r)
			beg := br.ntou4(f.r)
			var end uint64
			var seekfree uint64
			var nbytesfree uint32
			var nfree uint32
			var nbytesname uint32
			var units byte
			var compress uint32
			var seekinfo uint64
			var nbytesinfo uint32
			if ver < 1000000 { // small file 
				end = uint64(br.ntou4(f.r))
				seekfree = uint64(br.ntou4(f.r))
				nbytesfree = br.ntou4(f.r)
				nfree = uint32(br.ntoi4(f.r))
				nbytesname = br.ntou4(f.r)
				units = br.ntobyte(f.r)
				compress = br.ntou4(f.r)
				seekinfo = uint64(br.ntou4(f.r))
				nbytesinfo = br.ntou4(f.r)
			} else {
				end = br.ntou8(f.r)
				seekfree = br.ntou8(f.r)
				nbytesfree = br.ntou4(f.r)
				nfree = br.ntou4(f.r)
				nbytesname = br.ntou4(f.r)
				units = br.ntobyte(f.r)
				compress = br.ntou4(f.r)
				seekinfo = br.ntou8(f.r)
				nbytesinfo = br.ntou4(f.r)
			}
			println("ver:",ver)
			println("beg:",beg)
			println("end:",end)
			println("seekfree:",seekfree)
			println("nbytesfree:",nbytesfree)
			println("nfree:",nfree)
			println("nbytesname:",nbytesname)
			println("units:",units)
			println("compress:",compress)
			println("seekinfo:",seekinfo)
			println("nbytesinfo:",nbytesinfo)
			seekdir := beg
			println("seekdir:",seekdir)

			// --- read directory info -------
			nbytes := int64(nbytesname) + 22
			{
				f.seek(int64(nbytes), os.SEEK_SET)
				datimec := datime2time(br.ntou4(f.r))
				datimem := datime2time(br.ntou4(f.r))
				println("date-c:", datimec.String())
				println("date-m:", datimem.String())
				pos, err := f.r.Seek(0, os.SEEK_CUR)
				if err != nil {
					panic(err)
				}
				nbytes = pos
			}
			nbytes += 18 // fUUID.Sizeof()
			// assume the file may be above 2Gb if file version is > 4
			if ver >= 40000 {
				nbytes += 12
			}
			println("nbytes:",nbytes)
			{
				pos,_ = f.r.Seek(0, os.SEEK_CUR)
				println("current-pos:",pos)
			}
			f.seek(int64(beg), os.SEEK_SET)

			cbk3 := func(bufpos int64, f *FileReader) {
				println("--cbk3--")
				f.seek(bufpos, os.SEEK_SET)

				// pos is the buffer key location
				pos, err := f.r.Seek(0, os.SEEK_CUR)
				if err != nil {
					panic(err)
				}
				println("pos:",pos,bufpos,"->",nbytesname)
				f.r.Seek(int64(nbytesname), os.SEEK_SET)
				vers := br.ntou2(f.r)
				println("vers:",vers)
				datimec := datime2time(br.ntou4(f.r))
				datimem := datime2time(br.ntou4(f.r))
				println("date-c:",datimec.String())
				println("date-m:",datimem.String())
				_nbyteskeys := br.ntou4(f.r)
				_nbytesname := br.ntou4(f.r)
				println("nbyteskeys:",_nbyteskeys)
				println("nbytesname:",_nbytesname, nbytesname)
				
				if vers > 1000 {
					seekdir := br.ntou8(f.r)
					seekparent := br.ntou8(f.r)
					seekkeys := br.ntou8(f.r)
					println("seekdir:", seekdir)
					println("seekparent:",seekparent)
					println("seekkeys:",seekkeys)
				} else {
					seekdir := br.ntou4(f.r)
					seekparent := br.ntou4(f.r)
					seekkeys := br.ntou4(f.r)
					println("seekdir:", seekdir)
					println("seekparent:",seekparent)
					println("seekkeys:",seekkeys)
				}
				if (vers % 1000) > 1 {
					// skip uuid
					br.skip(f.r, 18)
				}

				// --- read TKey::FillBuffer info
				f.r.Seek(bufpos, os.SEEK_SET)
				//key := br.readKey(f.r)
				//println("keyvers:",key.version)
				nbytes := br.ntou4(f.r)
				keyver := br.ntoi2(f.r)
				println("nbytes:", nbytes)
				println("keyver:", keyver)

			}
			bufpos := nbytes
			if 300 > nbytes {
				bufpos = 300
			}
			cbk3(bufpos, f)
		}
		cbk2(4, f)
	}
	cbk1(0, f)

	return keys
}

func NewFileReader(name string) (f *FileReader, err error) {
	f = &FileReader{}
	f.url = name
	f.order = binary.BigEndian
	f.keys = []Key{}

	f.r, err = os.Open(name)
	if err != nil {
		return nil, err
	}

	br := f.breader()

	{
		hdr := make([]byte, 4)
		err = binary.Read(f.r, f.order, &hdr)
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
	println("end=", f.hdr.end)
	println("seekfree=",f.hdr.seekfree)
	println("nbytesfree=",f.hdr.nbytesfree)
	println("units=", f.hdr.units)
	println("nbytesinfo=", f.hdr.nbytesinfo)
	println("seekinfo:", f.hdr.seekinfo)
	_, err = f.r.Seek(int64(f.hdr.beg), os.SEEK_SET)
	if err != nil {
		return nil, err
	}

	key := br.readKey(f.r)
	f.keys = append(f.keys, key)
	pos, err := f.r.Seek(0, os.SEEK_CUR)
	println("current-offset:", pos)
	println("key-name:", f.keys[0].name)
	f.name = key.name
	f.title = key.title

	println("\n=== read-keys ===")
	f.keys = br.readKeys(f)
	println("=== read-keys === [done]")
	println("keys:",len(f.keys))
	return f, err
}

func (f *FileReader) breader() breader {
	return breader{f.order}
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

func (f *FileReader) GetKey(name string, cycle int) *Key {
	for i, _ := range f.keys {
		k := &f.keys[i]
		if k.name == name && int(k.cycle) == cycle {
			return k
		}
	}
	//FIXME: loop over directories too.

	return nil
}

func (f *FileReader) seek(offset int64, pos int) (ret int64, err error) {
	switch pos {
	case os.SEEK_SET:
		f.offset = uint64(offset) + f.archiveoffset
		return f.r.Seek(int64(f.offset), os.SEEK_SET)
	case os.SEEK_CUR:
		f.offset += uint64(offset)
		return f.r.Seek(int64(f.offset), os.SEEK_SET)
	case os.SEEK_END:
		if f.archiveoffset != 0 {
			panic("seek: seeking from end in archive is not (yet?) supported")
		}
		f.offset = f.hdr.end - uint64(offset)
		return f.r.Seek(int64(f.offset), os.SEEK_SET)

	default:
		err = fmt.Errorf("seek: unknown seek option (%d)", pos)
		panic(err)
	}
	panic("unreachable")
}

// EOF
