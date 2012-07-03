package groot

import (
	"reflect"
)

// leaf of ints

type LeafI struct {
	base baseLeaf
	min int32
	max int32
	data []int
}

func (leaf *LeafI) toBaseLeaf() *baseLeaf {
	return &leaf.base
}

func (leaf *LeafI) Class() Class {
	panic("not implemented")
}

func (leaf *LeafI) Name() string {
	return leaf.base.name
}

func (leaf *LeafI) Title() string {
	return leaf.base.title
}

func (leaf *LeafI) ROOTDecode(b *Buffer) (err error) {

	vers, pos, bcnt := b.read_version()
	dprintf("leafI-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]int, int(leaf.base.length))
	dprintf("leafI min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	return
}

func (leaf *LeafI) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

// leaf of ints-64

type LeafL struct {
	base baseLeaf
	min int32
	max int32
	data []int64
}

func (leaf *LeafL) toBaseLeaf() *baseLeaf {
	return &leaf.base
}

func (leaf *LeafL) Class() Class {
	panic("not implemented")
}

func (leaf *LeafL) Name() string {
	return leaf.base.name
}

func (leaf *LeafL) Title() string {
	return leaf.base.title
}

func (leaf *LeafL) ROOTDecode(b *Buffer) (err error) {

	vers, pos, bcnt := b.read_version()
	dprintf("leafL-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]int64, int(leaf.base.length))
	dprintf("leafL min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	return
}

func (leaf *LeafL) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

// leaf of floats

type LeafF struct {
	base baseLeaf
	min int32
	max int32
	data []float32
}

func (leaf *LeafF) toBaseLeaf() *baseLeaf {
	return &leaf.base
}

func (leaf *LeafF) Class() Class {
	panic("not implemented")
}

func (leaf *LeafF) Name() string {
	return leaf.base.name
}

func (leaf *LeafF) Title() string {
	return leaf.base.title
}

func (leaf *LeafF) ROOTDecode(b *Buffer) (err error) {

	vers, pos, bcnt := b.read_version()
	dprintf("leafF-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]float32, int(leaf.base.length))
	dprintf("leafF min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	return
}

func (leaf *LeafF) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

// leaf of doubles

type LeafD struct {
	base baseLeaf
	min int32
	max int32
	data []float64
}

func (leaf *LeafD) toBaseLeaf() *baseLeaf {
	return &leaf.base
}

func (leaf *LeafD) Class() Class {
	panic("not implemented")
}

func (leaf *LeafD) Name() string {
	return leaf.base.name
}

func (leaf *LeafD) Title() string {
	return leaf.base.title
}

func (leaf *LeafD) ROOTDecode(b *Buffer) (err error) {

	vers, pos, bcnt := b.read_version()
	dprintf("leafD-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]float64, int(leaf.base.length))
	dprintf("leafD min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	return
}

func (leaf *LeafD) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

func init() {
	{
		f := func() reflect.Value {
			o := &LeafI{}
			return reflect.ValueOf(o)
		}
		Factory.db["TLeafI"] = f
		Factory.db["*groot.LeafI"] = f
	}

	{
		f := func() reflect.Value {
			o := &LeafL{}
			return reflect.ValueOf(o)
		}
		Factory.db["TLeafL"] = f
		Factory.db["*groot.LeafL"] = f
	}

	{
		f := func() reflect.Value {
			o := &LeafF{}
			return reflect.ValueOf(o)
		}
		Factory.db["TLeafF"] = f
		Factory.db["*groot.LeafF"] = f
	}

	{
		f := func() reflect.Value {
			o := &LeafD{}
			return reflect.ValueOf(o)
		}
		Factory.db["TLeafD"] = f
		Factory.db["*groot.LeafD"] = f
	}
}

// EOF
