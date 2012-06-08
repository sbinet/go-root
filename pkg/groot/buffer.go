package groot

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Buffer struct {
	order binary.ByteOrder // byte order of underlying data source
	data  []byte           // data source
	buf   *bytes.Buffer    // buffer for more efficient i/o from r
	klen  uint32           // to compute refs (used in read_class, read_object)
}

func NewBuffer(data []byte, order binary.ByteOrder, klen uint32) (b *Buffer, err error) {
	b = &Buffer{
		order: order,
		data:  data,
		klen:  klen,
	}
	b.buf = bytes.NewBuffer(b.data[:])
	return
}

func NewBufferFromKey(k *Key) (b *Buffer, err error) {
	buf, err := k.Buffer()
	if err != nil {
		return
	}
	return NewBuffer(buf, k.file.order, uint32(k.keysz))
}

func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *Buffer) clone() *Buffer {
	bb, err := NewBuffer(b.buf.Bytes(), b.order, b.klen)
	if err != nil {
		return nil
	}
	return bb
}

func (b *Buffer) read_nbytes(nbytes int) (o []byte) {
	o = make([]byte, nbytes)
	err := binary.Read(b.buf, b.order, o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntoi2() (o int16) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntoi4() (o int32) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntoi8() (o int64) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntobyte() (o byte) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntou2() (o uint16) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntou4() (o uint32) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntou8() (o uint64) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntof() (o float32) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntod() (o float64) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) read_array_F() (o []float32) {
	n := int(b.ntou4())
	o = make([]float32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntof()
	}
	return
}

func (b *Buffer) read_array_D() (o []float64) {
	n := int(b.ntou4())
	o = make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntod()
	}
	return
}

func (b *Buffer) read_array_S() (o []int16) {
	n := int(b.ntou4())
	o = make([]int16, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi2()
	}
	return
}

func (b *Buffer) read_array_I() (o []int32) {
	n := int(b.ntou4())
	o = make([]int32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi4()
	}
	return
}

func (b *Buffer) read_array_L() (o []int64) {
	n := int(b.ntou4())
	o = make([]int64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi8()
	}
	return
}

func (b *Buffer) read_array_C() (o []byte) {
	n := int(b.ntou4())
	o = make([]byte, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntobyte()
	}
	return
}

func (b *Buffer) read_static_array() (o []uint32) {
	n := int(b.ntou4())
	o = make([]uint32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou4()
	}
	return
}

func (b *Buffer) read_fast_array_F(n int) (o []float32) {
	o = make([]float32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntof()
	}
	return
}

func (b *Buffer) read_fast_array_D(n int) (o []float64) {
	o = make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntod()
	}
	return
}

func (b *Buffer) read_fast_array_S(n int) (o []int16) {
	o = make([]int16, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi2()
	}
	return
}

func (b *Buffer) read_fast_array_I(n int) (o []int32) {
	o = make([]int32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi4()
	}
	return
}

func (b *Buffer) read_fast_array_L(n int) (o []int64) {
	o = make([]int64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi8()
	}
	return
}

func (b *Buffer) read_fast_array_C(n int) (o []byte) {
	o = make([]byte, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntobyte()
	}
	return
}

func (b *Buffer) read_fast_array_tstring(n int) (o []string) {
	o = make([]string, n)
	for i := 0; i < n; i++ {
		o[i] = b.read_tstring()
	}
	return
}

func (b *Buffer) read_fast_array(n int) (o []uint32) {
	o = make([]uint32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou4()
	}
	return
}

func (b *Buffer) read_tstring() string {
	n := int(b.ntobyte())
	if n == 255 {
		// large string
		n = int(b.ntou4())
	}
	if n == 0 {
		return ""
	}
	v := b.ntobyte()
	if v == 0 {
		return ""
	}
	o := make([]byte, n)
	o[0] = v
	err := binary.Read(b.buf, b.order, o[1:])
	if err != nil {
		panic(err)
	}
	return string(o)
}

//FIXME
// readBasicPointerElem
// readBasicPointer

func (b *Buffer) read_string(max int) string {
	o := []byte{}
	n := 0
	var v byte
	for {
		v = b.ntobyte()
		if v == 0 {
			break
		}
		n += 1
		if max > 0 && n >= max {
			break
		}
		o = append(o, v)
	}
	return string(o)
}

func (b *Buffer) read_std_string() string {
	nwh := b.ntobyte()
	nchars := int32(nwh)
	if nwh == 255 {
		nchars = b.ntoi4()
	}
	if nchars < 0 {
		panic("groot.readStdString: negative char number")
	}
	return b.read_string(int(nchars))
}

func (b *Buffer) read_version() (vers uint16, pos, bcnt uint32) {

	bcnt = b.ntou4()
	if (int64(bcnt) & ^kByteCountMask) != 0 {
		bcnt = uint32(int64(bcnt) & ^kByteCountMask)
	} else {
		panic("groot.Buffer.read_version: too old file")
	}
	vers = b.ntou2()
	return
}

func (b *Buffer) read_object() (o Object) {
	clsname, bcnt, isref := b.read_class()
	printf(">>[%s] [%v] [%v]\n", clsname, bcnt, isref)
	if isref {
		obj_offset := bcnt - kMapOffset - b.klen
		bb := b.clone()
		bb.read_nbytes(int(obj_offset))
		ii := bb.ntou4()
		if (ii & kByteCountMask) != 0 {
			scls := bb.read_class_tag()
			if scls == "" {
				panic("groot.Buffer.read_object: read_class_tag did not find a class name")
			}
		} else {
			/* boo */
		}
		/*
		 // in principle at this point m_pos should be
		 //   m_buffer+startpos+sizeof(unsigned int)
		 // but enforce it anyway : 
		 m_pos = m_buffer+startpos+sizeof(unsigned int); 
		*/
	} else {
		if clsname == "" {
			o = nil
		} else {

			factory := Factory.Get(clsname)
			if factory == nil {
				dprintf("**err** no factory for class [%s]\n", clsname)
				return
			}

			vv := factory()
			o = vv.Interface().(Object)
			if vv, ok := vv.Interface().(ROOTStreamer); ok {
				err := vv.ROOTDecode(b.clone())
				if err != nil {
					panic(err)
				}
			} else {
				dprintf("**err** class [%s] does not satisfy the ROOTStreamer interface\n", clsname)
			}
		}
	}
	return o
}

func (b *Buffer) read_class() (name string, bcnt uint32, isref bool) {

	//var bufvers = 0
	i := b.ntou4()

	if i == kNullTag {
		/*empty*/
	} else if (i & kByteCountMask) != 0 {
		//bufvers = 1
		clstag := b.read_class_tag()
		if clstag == "" {
			panic("groot.breader.readClass: empty class tag")
		}
		name = clstag
		bcnt = uint32(int64(i) & ^kByteCountMask)
	} else {
		bcnt = uint32(i)
		isref = true
	}
	printf("--[%s] [%v] [%v]\n", name, bcnt, isref)
	return
}

func (b *Buffer) read_class_tag() (clstag string) {
	tag := b.ntou4()

	if tag == kNewClassTag {
		clstag = b.read_string(80)
		printf("--class+tag: [%v]\n", clstag)
	} else if (tag & kClassMask) != 0 {
		clstag = b.clone().read_class_tag()
		printf("--class-tag: [%v]\n", clstag)
	} else {
		panic(fmt.Errorf("groot.read_class_tag: unknown class-tag [%v]", tag))
	}
	return
}

func (b *Buffer) read_tnamed() (name, title string) {

	vers, pos, bcnt := b.read_version()
	id := b.ntou4()
	bits := b.ntou4()
	bits |= kIsOnHeap // by definition de-serialized object is on heap
	if (bits & kIsReferenced) == 0 {
		_ = b.read_nbytes(2)
	}
	name = b.read_tstring()
	title = b.read_tstring()
	printf("read_tnamed: vers=%v pos=%v bcnt=%v id=%v bits=%v name='%v' title='%v'\n",
		vers, pos, bcnt, id, bits, name, title)
	//FIXME: buffer.check_byte_count(pos,bcnt,"TNamed")

	return
}

func (b *Buffer) read_elements() (elmts []Object) {
	name, bcnt, isref := b.read_class()
	printf("read_elements: name='%v' bcnt=%v isref=%v\n",
		name, bcnt, isref)
	elmts = b.read_obj_array()
	return elmts
}

func (b *Buffer) read_obj_array() (elmts []Object) {

	vers, pos, bcnt := b.read_version()
	if vers > 2 {
		// skip version
		b.read_nbytes(2)
		// skip object bits and unique id
		b.read_nbytes(8)
	}
	name := "??"
	if vers > 1 {
		name = b.read_tstring()
	}
	title := b.read_tstring()

	nobjs := int(b.ntoi4())
	lbound := b.ntoi4()

	printf("read_obj_array: vers=%v pos=%v bcnt=%v name='%v' title='%v' nobjs=%v lbound=%v\n",
		vers, pos, bcnt, name, title, nobjs, lbound)

	elmts = make([]Object, nobjs)
	for i := 0; i < nobjs; i++ {
		obj := b.read_object()
		elmts[i] = obj
	}
	//FIXME: buffer.check_byte_count(s,c,"TObjArray")
	return elmts
}

/*
    short v;
    unsigned int s, c;
    if(!a_buffer.read_version(v,s,c)) return false;

    //::printf("debug : ObjArray::stream : version %d count %d\n",v,c);

   {uint32 id,bits;
    if(!Object_stream(a_buffer,id,bits)) return false;}
    std::string name;
    if(!a_buffer.read(name)) return false;
    int nobjects;
    if(!a_buffer.read(nobjects)) return false;
    int lowerBound;
    if(!a_buffer.read(lowerBound)) return false;

    //::printf("debug : ObjArray : nobject \"%s\" %d %d\n",
    //  name.c_str(),nobjects,lowerBound);

    for (int i=0;i<nobjects;i++) {
      //::printf("debug : ObjArray :    n=%d i=%d ...\n",nobjects,i);
      iro* obj;
      if(!a_buffer.read_object(a_fac,a_args,obj)){
        a_buffer.out() << "inlib::rroot::ObjArray::stream :"
                       << " can't read object."
                       << std::endl;
        return false;
      }
      //::printf("debug : ObjArray :    n=%d i=%d : ok\n",nobjects,i);
      if(obj) {
        T* to = inlib::cast<iro,T>(*obj);
        if(!to) {
          a_buffer.out() << "inlib::rroot::ObjArray::stream :"
                         << " inlib::cast failed."
                         << std::endl;
        } else {
          push_back(to);
        }
      } else {
        //a_accept_null for branch::stream m_baskets.
        if(a_accept_null) this->push_back(0);
      }
    }

    return a_buffer.check_byte_count(s,c,"TObjArray");
  }
*/

//FIXME
// readObjectAny
// readTList
// readTObjArray
// readTClonesArray
// readTCollection
// readTHashList
// readTNamed
// readTCanvas

// EOF
