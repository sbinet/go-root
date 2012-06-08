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

func (b *Buffer) readArrayF() (o []float32) {
	n := int(b.ntou4())
	o = make([]float32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntof()
	}
	return
}

func (b *Buffer) readArrayD() (o []float64) {
	n := int(b.ntou4())
	o = make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntod()
	}
	return
}

func (b *Buffer) readArrayS() (o []int16) {
	n := int(b.ntou4())
	o = make([]int16, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi2()
	}
	return
}

func (b *Buffer) readArrayI() (o []int32) {
	n := int(b.ntou4())
	o = make([]int32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi4()
	}
	return
}

func (b *Buffer) readArrayL() (o []int64) {
	n := int(b.ntou4())
	o = make([]int64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi8()
	}
	return
}

func (b *Buffer) readArrayC() (o []byte) {
	n := int(b.ntou4())
	o = make([]byte, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntobyte()
	}
	return
}

func (b *Buffer) readStaticArray() (o []uint32) {
	n := int(b.ntou4())
	o = make([]uint32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou4()
	}
	return
}

func (b *Buffer) readFastArrayF(n int) (o []float32) {
	o = make([]float32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntof()
	}
	return
}

func (b *Buffer) readFastArrayD(n int) (o []float64) {
	o = make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntod()
	}
	return
}

func (b *Buffer) readFastArrayS(n int) (o []int16) {
	o = make([]int16, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi2()
	}
	return
}

func (b *Buffer) readFastArrayI(n int) (o []int32) {
	o = make([]int32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi4()
	}
	return
}

func (b *Buffer) readFastArrayL(n int) (o []int64) {
	o = make([]int64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi8()
	}
	return
}

func (b *Buffer) readFastArrayC(n int) (o []byte) {
	o = make([]byte, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntobyte()
	}
	return
}

func (b *Buffer) readFastArrayTString(n int) (o []string) {
	o = make([]string, n)
	for i := 0; i < n; i++ {
		o[i] = b.readTString()
	}
	return
}

func (b *Buffer) readFastArray(n int) (o []uint32) {
	o = make([]uint32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou4()
	}
	return
}

func (b *Buffer) readTString() string {
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

func (b *Buffer) readString(max int) string {
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

func (b *Buffer) readStdString() string {
	nwh := b.ntobyte()
	nchars := int32(nwh)
	if nwh == 255 {
		nchars = b.ntoi4()
	}
	if nchars < 0 {
		panic("groot.readStdString: negative char number")
	}
	return b.readString(int(nchars))
}

func (b *Buffer) readVersion() (vers uint16, pos, bcnt uint32) {

	bcnt = b.ntou4()
	if (int64(bcnt) & ^kByteCountMask) != 0 {
		bcnt = uint32(int64(bcnt) & ^kByteCountMask)
	} else {
		panic("groot.breader.readVersion: too old file")
	}
	vers = b.ntou2()
	return
}

func (b *Buffer) readObject() (o Object) {
	clsname, bcnt, isref := b.readClass()
	dprintf(">>[%s] [%v] [%v]\n", clsname, bcnt, isref)
	if isref {

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
			// if vv,ok := vv.Interface().(FileSetter); ok {
			// 	err := vv.SetFile(k.file)
			// 	if err != nil {
			// 		return v
			// 	}
			// }
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

func (b *Buffer) readClass() (name string, bcnt uint32, isref bool) {

	//var bufvers = 0
	i := b.ntou4()

	if i == kNullTag {
		/*empty*/
	} else if (i & kByteCountMask) != 0 {
		//bufvers = 1
		clstag := b.readClassTag()
		if clstag == "" {
			panic("groot.breader.readClass: empty class tag")
		}
		name = clstag
		bcnt = uint32(int64(i) & ^kByteCountMask)
	} else {
		bcnt = uint32(i)
		isref = true
	}
	dprintf("--[%s] [%v] [%v]\n", name, bcnt, isref)
	return
}

func (b *Buffer) readClassTag() (clstag string) {
	tag := b.ntou4()

	if tag == kNewClassTag {
		clstag = b.readString(80)
		dprintf("--class+tag: [%v]\n", clstag)
	} else if (tag & kClassMask) != 0 {
		clstag = b.clone().readClassTag()
		dprintf("--class-tag: [%v]\n", clstag)
	} else {
		panic(fmt.Errorf("groot.readClassTag: unknown class-tag [%v]", tag))
	}
	return
}

// func (b *Buffer) readObject(r io.Reader) (id, bits uint32) {
// 	/*v,pos,bcnt := */ br.readVersion(r)
// 	id = br.ntou8(r)
// 	bits = br.ntou8(r)
// 	return
// }

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
