package groot

import (
	"reflect"
)

type Tree struct {
	name     string
	title    string
	entries  uint64
	branches []Branch
}

func (tree *Tree) Class() Class {
	panic("not implemented")
}

func (tree *Tree) Name() string {
	return tree.name
}

func (tree *Tree) Title() string {
	return tree.title
}

func (tree *Tree) ROOTDecode(b *Buffer) (err error) {

	vers, pos, bcnt := b.read_version()
	dprintf("vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	tree.name, tree.title = b.read_tnamed()
	dprintf("name='%v' title='%v'\n", tree.name, tree.title)

	return
}

func (tree *Tree) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

func init() {
	f := func() reflect.Value {
		o := &Tree{branches: make([]Branch, 0)}
		return reflect.ValueOf(o)
	}
	Factory.db["TTree"] = f
	Factory.db["*groot.Tree"] = f
}

// check interfaces
var _ Object = (*Tree)(nil)
var _ ROOTStreamer = (*Tree)(nil)

// EOF
