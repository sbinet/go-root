package groot

import (
	"fmt"
	"reflect"
)

type Tree struct {
	file     *File
	name     string
	title    string
	entries  uint64
	branches []Branch
}

func NewTree(file *File, name, title string) (tree *Tree, err error) {
	tree = &Tree{
	file: file,
	name: name,
	title: title,
	entries: 0,
	branches: make([]Branch, 0),
	}
	return
}

func (tree *Tree) SetFile(f *File) (err error) {
	if tree.file != nil {
		err = fmt.Errorf("groot: cannot migrate a Tree to another file")
		return
	}
	tree.file = f
	return
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
	printf("vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	tree.name, tree.title = b.read_tnamed()
	printf("name='%v' title='%v'\n", tree.name, tree.title)

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
var _ FileSetter = (*Tree)(nil)

// EOF
