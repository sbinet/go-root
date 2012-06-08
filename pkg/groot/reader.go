package groot

import (
	"encoding/binary"
	"fmt"
	"os"
)

type unzip_fct func()

type File struct {
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

func NewFileReader(name string) (f *File, err error) {
	f = &File{
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

func (f *File) seek(offset int64, pos int) (ret int64, err error) {
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

func (f *File) initialize() (err error) {
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

	buf := make([]byte, int(nbytes))

	// read directory info
	_, err = f.f.ReadAt(buf, f.beg)
	if err != nil {
		return err
	}

	b, err := NewBuffer(buf[f.nbytes_name:], f.order, 0)
	if err != nil {
		return err
	}

	err = f.root_dir.from_buffer(b)
	if err != nil {
		return err
	}

	nk := sz_int32     // Key::fNumberOfBytes
	nk += sz_int16     // Key::fVersion
	nk += 2 * sz_int32 // Key::fObjectSize, Date
	nk += 2 * sz_int16 // Key::fKeyLength, fCycle
	nk += 2 * sz_int32 // Key::fSeekKey, fSeekParentDirectory
	// WARNING: the above is sz_int32 since we are at beginning of file

	b, err = NewBuffer(buf[nk:], f.order, 0)
	if err != nil {
		return err
	}

	printf("nk: %v\n", nk)
	printf("buf: %v\n", buf[:])
	printf("rst: %v\n", buf[nk:])
	printf("TFile: %v\n", []byte{'T', 'F', 'i', 'l', 'e'})
	cname := b.readTString()
	if cname != "TFile" {
		return fmt.Errorf("groot: expected [TFile]. got [%v]", cname)
	}
	printf("f-clsname [%v]\n", cname)

	cname = b.readTString()
	printf("f-cname   [%v]\n", cname)

	f.title = b.readTString()
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

func (f *File) read_header() (err error) {
	buf := make([]byte, 64)
	_, err = f.f.ReadAt(buf, 0)
	if err != nil {
		return err
	}

	b, err := NewBuffer(buf, f.order, 0)
	{
		hdr := b.read_nbytes(4)
		printf("hdr: %v\n", string(hdr))
		if string(hdr) != "root" {
			return fmt.Errorf(
				"groot: file [%s] is not a ROOT file (%v)",
				f.name, string(hdr))
		}
	}
	f.version = b.ntou4()
	f.beg = int64(b.ntou4())
	printf("beg: %v\n", f.beg)
	if f.version >= 1000000 {
		f.end = int64(b.ntou8())
		f.seek_free = int64(b.ntou8())
	} else {
		f.end = int64(b.ntou4())
		f.seek_free = int64(b.ntou4())
	}
	printf("end: %v\n", f.end)
	printf("seek-free: %v\n", f.seek_free)
	f.nbytes_free = b.ntou4()
	/*nfree*/ b.ntoi4()
	f.nbytes_name = b.ntou4()
	printf("nbytes-free: %v\n", f.nbytes_free)
	printf("nbytes-name: %v\n", f.nbytes_name)
	/*units*/ b.ntobyte()
	/*compress*/ b.ntou4()
	if f.version >= 1000000 {
		f.seek_info = int64(b.ntou8())
	} else {
		f.seek_info = int64(b.ntou4())
	}
	f.nbytes_info = b.ntou4()
	printf("seek-info: %v\n", f.seek_info)
	printf("nbytes-info: %v\n", f.nbytes_info)

	// read streamer infos
	return f.read_streamer_infos()
}

func (f *File) read_streamer_infos() (err error) {

	if true {
		return
	}
	lst := make(List, 0)
	dprintf("lst: %v\n", lst)

	var buf []byte

	if f.seek_info > 0 && f.seek_info < f.end {
		buf = make([]byte, int(f.nbytes_info))

		_, err = f.f.ReadAt(buf, f.seek_info)
		if err != nil {
			return err
		}
		key, err := NewKey(f, 0, 0)
		if err != nil {
			return err
		}
		b, err := NewBuffer(buf, f.order, 0)
		err = key.init_from_buffer(b)
		if err != nil {
			return err
		}
		obj, ok := key.Value().(List)
		if ok {
			lst = obj
		}
	}
	dprintf("buf: %v\n", len(buf))
	
	return err
}

/*
void TFile::ReadStreamerInfo()
{
// Read the list of StreamerInfo from this file
// The key with name holding the list of TStreamerInfo objects is read.
// The corresponding TClass objects are updated.

   TList *list = 0;
   if (fSeekInfo > 0 && fSeekInfo < fEND) {
      TKey *key = new TKey();
      char *buffer = new char[fNbytesInfo+1];
      char *buf    = buffer;
      Seek(fSeekInfo);
      ReadBuffer(buf,fNbytesInfo);
      key->ReadBuffer(buf);
      TFile *filesave = gFile;
      gFile = this; // used in ReadObj
      list = (TList*)key->ReadObj();
      if (!list) {
         gDirectory->GetListOfKeys()->Remove(key);
         MakeZombie();
      }
      gFile = filesave;
      delete [] buffer;
      delete key;
   } else {
      list = (TList*)Get("StreamerInfo"); //for versions 2.26 (never released)
   }

   if (list == 0) return;
   if (gDebug > 0) Info("ReadStreamerInfo", "called for file %s",GetName());

   // loop on all TStreamerInfo classes
   TStreamerInfo *info;
   TIter next(list);
   while ((info = (TStreamerInfo*)next())) {
      if (info->IsA() != TStreamerInfo::Class()) {
         Warning("ReadStreamerInfo","%s: not a TStreamerInfo object", GetName());
         continue;
      }
      info->BuildCheck();
      Int_t uid = info->GetNumber();
      Int_t asize = fClassIndex->GetSize();
      if (uid >= asize && uid <100000) fClassIndex->Set(2*asize);
      if (uid >= 0 && uid < fClassIndex->GetSize()) fClassIndex->fArray[uid] = 1;
      else {
         printf("ReadStreamerInfo, class:%s, illegal uid=%d\n",info->GetName(),uid);
      }
      if (gDebug > 0) printf(" -class: %s version: %d info read at slot %d\n",info->GetName(), info->GetClassVersion(),uid);
   }
   fClassIndex->fArray[0] = 0;
   list->Clear();  //this will delete all TStreamerInfo objects with kCanDelete bit set
   delete list;
}

*/

func (f *File) Name() string {
	return f.name
}

func (f *File) Version() uint32 {
	return f.version
}

func (f *File) Dir() *Directory {
	return &f.root_dir
}

func (f *File) ByteOrder() binary.ByteOrder {
	return f.order
}

// EOF
