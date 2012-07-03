package groot

import (
	"reflect"
)

// leaf of bytes

type LeafB struct {
	base baseLeaf
	min int32
	max int32
	data []byte
}

func (leaf *LeafB) toBaseLeaf() *baseLeaf {
	return &leaf.base
}

func (leaf *LeafB) Class() Class {
	panic("not implemented")
}

func (leaf *LeafB) Name() string {
	return leaf.base.name
}

func (leaf *LeafB) Title() string {
	return leaf.base.title
}

func (leaf *LeafB) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("[leafB] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]byte, int(leaf.base.length))
	printf("leafI min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	b.check_byte_count(pos,bcnt,spos, "LeafB")
	return
}

func (leaf *LeafB) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

// leaf of shorts

type LeafS struct {
	base baseLeaf
	min int32
	max int32
	data []int8
}

func (leaf *LeafS) toBaseLeaf() *baseLeaf {
	return &leaf.base
}

func (leaf *LeafS) Class() Class {
	panic("not implemented")
}

func (leaf *LeafS) Name() string {
	return leaf.base.name
}

func (leaf *LeafS) Title() string {
	return leaf.base.title
}

func (leaf *LeafS) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("[leafS] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]int8, int(leaf.base.length))
	printf("leafI min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	b.check_byte_count(pos,bcnt,spos, "LeafS")
	return
}

func (leaf *LeafS) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

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
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("leafI-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]int, int(leaf.base.length))
	printf("leafI min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	b.check_byte_count(pos,bcnt,spos, "LeafI")
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
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("leafL-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]int64, int(leaf.base.length))
	printf("leafL min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	b.check_byte_count(pos,bcnt,spos, "LeafL")
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
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("leafF-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]float32, int(leaf.base.length))
	printf("leafF min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	b.check_byte_count(pos,bcnt,spos, "LeafF")
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
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("leafD-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data = make([]float64, int(leaf.base.length))
	printf("leafD min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	b.check_byte_count(pos,bcnt,spos, "LeafD")
	return
}

func (leaf *LeafD) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

// leaf of a string

type LeafC struct {
	base baseLeaf
	min int32
	max int32
	data string
}

func (leaf *LeafC) toBaseLeaf() *baseLeaf {
	return &leaf.base
}

func (leaf *LeafC) Class() Class {
	panic("not implemented")
}

func (leaf *LeafC) Name() string {
	return leaf.base.name
}

func (leaf *LeafC) Title() string {
	return leaf.base.title
}

func (leaf *LeafC) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("[leafC] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	printf("leafC min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	b.check_byte_count(pos,bcnt,spos, "LeafC")
	return
}

func (leaf *LeafC) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

// leaf of bool

type LeafO struct {
	base baseLeaf
	min int32
	max int32
	data []bool
}

func (leaf *LeafO) toBaseLeaf() *baseLeaf {
	return &leaf.base
}

func (leaf *LeafO) Class() Class {
	panic("not implemented")
}

func (leaf *LeafO) Name() string {
	return leaf.base.name
}

func (leaf *LeafO) Title() string {
	return leaf.base.title
}

func (leaf *LeafO) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("[leafO] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	err = leaf.base.ROOTDecode(b)
	if err != nil {
		return err
	}
	leaf.min = b.ntoi4()
	leaf.max = b.ntoi4()
	leaf.data= make([]bool, int(leaf.base.length))
	printf("leafO min=%v max=%v len=%d\n", leaf.min, leaf.max, len(leaf.data))
	b.check_byte_count(pos,bcnt,spos, "LeafO")
	return
}

func (leaf *LeafO) ROOTEncode(b *Buffer) (err error) {
	//FIXME
	panic("not implemented")
}

func init() {

	{
		f := func() reflect.Value {
			o := &LeafO{}
			return reflect.ValueOf(o)
		}
		Factory.db["TLeafO"] = f
		Factory.db["*groot.LeafO"] = f
	}

	{
		f := func() reflect.Value {
			o := &LeafB{}
			return reflect.ValueOf(o)
		}
		Factory.db["TLeafB"] = f
		Factory.db["*groot.LeafB"] = f
	}

	{
		f := func() reflect.Value {
			o := &LeafS{}
			return reflect.ValueOf(o)
		}
		Factory.db["TLeafS"] = f
		Factory.db["*groot.LeafS"] = f
	}

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

	{
		f := func() reflect.Value {
			o := &LeafC{}
			return reflect.ValueOf(o)
		}
		Factory.db["TLeafC"] = f
		Factory.db["*groot.LeafC"] = f
	}
}

// check interfaces
var _ Object = (*LeafO)(nil)
var _ ROOTStreamer = (*LeafO)(nil)

var _ Object = (*LeafB)(nil)
var _ ROOTStreamer = (*LeafB)(nil)

var _ Object = (*LeafS)(nil)
var _ ROOTStreamer = (*LeafS)(nil)

var _ Object = (*LeafI)(nil)
var _ ROOTStreamer = (*LeafI)(nil)

var _ Object = (*LeafL)(nil)
var _ ROOTStreamer = (*LeafL)(nil)

var _ Object = (*LeafF)(nil)
var _ ROOTStreamer = (*LeafF)(nil)

var _ Object = (*LeafD)(nil)
var _ ROOTStreamer = (*LeafD)(nil)

var _ Object = (*LeafC)(nil)
var _ ROOTStreamer = (*LeafC)(nil)

// EOF
