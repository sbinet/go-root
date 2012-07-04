package groot

import (
	"fmt"
	"reflect"
)

type List struct {
	elmts []Object
}

func (lst *List) Class() string {
	return "TList"
}

func (lst *List) Name() string {
	return "list-name"
}

func (lst *List) Title() string {
	return "list-title"
}

func (lst *List) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()
	lst.elmts = make([]Object, 0)

	vers, pos, bcnt := b.read_version()
	printf("[list] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)

	if vers <= 3 {
		err = fmt.Errorf("sorry, too old version of TList (%d)", vers)
		return
	}
	id := b.ntou4()
	bits := b.ntou4()
	bits |= kIsOnHeap // by definition de-serialized object is on heap
	if (bits & kIsReferenced) == 0 {
		_ = b.read_nbytes(2)
	}

	name := b.read_tstring()
	nobjs := int(b.ntoi4())
	lst.elmts = make([]Object, 0, nobjs)

	printf("id=%v bits=%v\n", id, bits)
	printf("name=%v nobjs=%v\n", name, nobjs)

	for i := 0; i < nobjs; i++ {
		printf("---> %v/%v\n", i+1, nobjs)
		obj := b.read_object()
		opt := b.read_tstring()

		// obj can be nil if the class had a custom streamer and we do not
		// have the shared library nor a streamerinfo.
		if obj != nil {
			lst.elmts = append(lst.elmts, obj)
		}
		printf("---> %v/%v == opt='%s' obj=%v\n", i+1, nobjs, opt, obj)
	}

	b.check_byte_count(pos, bcnt, spos, "TList")
	return err
}

func (lst *List) ROOTEncode(b *Buffer) (err error) {
	panic("groot.List.ROOTEncode not implemented")
}

/*
  bool stream(buffer& a_buffer){
    // Stream all objects in the collection to or from the I/O buffer.
    _clear();

    short v;
    unsigned int s, c;
    if(!a_buffer.read_version(v,s,c)) return false;
   {uint32 id,bits;
    if(!Object_stream(a_buffer,id,bits)) return false;}
    if(!a_buffer.read(m_name)) return false;

    int nobjects;
    if(!a_buffer.read(nobjects)) return false;

    for (int i = 0; i < nobjects; i++) {
      dummy_fac fac(a_buffer.out());
      ifac::args args;
      iro* obj;
      if(!a_buffer.read_object(fac,args,obj)) {
        a_buffer.out() << "inlib::rroot::List::stream :"
                       << " can't read object."
                       << " index " << i
                       << " over " << nobjects << " objects."
                       << std::endl;
        _clear();
        return false;
      }
      unsigned char nch;
      if(!a_buffer.read(nch)) {
        _clear();
        return false;
      }
      if(nch) {
        char readOption[256];
        if(!a_buffer.read_fast_array(readOption,nch)) {
          _clear();
          return false;
        }
        readOption[nch] = 0;
      } else {
      }
      if(obj) push_back(obj);
    }

    if(!a_buffer.check_byte_count(s,c,"TList")) {
      _clear();
      return false;
    }
    return true;
  }
*/

func init() {

	//FIXME...
	 new_lst := func() reflect.Value {
	 	o := List{}
	 	return reflect.ValueOf(&o)
	 }
	Factory.db["TList"] = new_lst
	Factory.db["*groot.List"] = new_lst
}

// check interfaces
var _ Object = (*List)(nil)
var _ ROOTStreamer = (*List)(nil)

// EOF
