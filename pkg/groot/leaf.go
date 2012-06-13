package groot

import (
	"reflect"
)

type LeafI struct {
	base baseLeaf
	min int32
	max int32
	data []int
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

func init() {
	f := func() reflect.Value {
		o := &LeafI{}
		return reflect.ValueOf(o)
	}
	Factory.db["TLeafI"] = f
	Factory.db["*groot.LeafI"] = f
}

// EOF
