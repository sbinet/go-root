package groot

import (
	"reflect"
)

type LeafElement struct {
	base   baseLeaf
	id     int // element serial number in fInfo
	ltype  int // leaf type
}

func (le *LeafElement) toBaseLeaf() *baseLeaf {
	return &le.base
}

func (le *LeafElement) Class() Class {
	panic("not implemented")
}

func (le *LeafElement) Name() string {
	return le.base.name
}

func (le *LeafElement) Title() string {
	return le.base.title
}

func (le *LeafElement) ROOTDecode(b *Buffer) (err error) {

	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	dprintf("vers=%v spos=%v pos=%v bcnt=%v\n", vers, spos, pos, bcnt)
	err = le.base.ROOTDecode(b)
	if err != nil {
		return err
	}

	le.id = int(b.ntoi4())
	le.ltype = int(b.ntoi4())

	b.check_byte_count(pos, bcnt, spos, "TLeafElement")
	return
}

func (le *LeafElement) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

func init() {
	f := func() reflect.Value {
		o := &LeafElement{}
		return reflect.ValueOf(o)
	}
	Factory.db["TLeafElement"] = f
	Factory.db["*groot.LeafElement"] = f
}

// check interfaces
var _ Object = (*LeafElement)(nil)
var _ ROOTStreamer = (*LeafElement)(nil)

