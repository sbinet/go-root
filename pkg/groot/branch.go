package groot

import (
	"reflect"
)

type Branch struct {
	name  string
	title string

	file *File
}

func (branch *Branch) Class() Class {
	panic("not implemented")
}

func (branch *Branch) Name() string {
	return branch.name
}

func (branch *Branch) Title() string {
	return branch.title
}

func (branch *Branch) ROOTDecode(b *Buffer) (err error) {

	vers, pos, bcnt := b.read_version()
	dprintf("vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	branch.name, branch.title = b.read_tnamed()
	dprintf("name='%v' title='%v'\n", branch.name, branch.title)

	return
}

func (branch *Branch) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

func init() {
	f := func() reflect.Value {
		o := &Branch{}
		return reflect.ValueOf(o)
	}
	Factory.db["TBranch"] = f
	Factory.db["*groot.Branch"] = f
}

// check interfaces
var _ Object = (*Branch)(nil)
var _ ROOTStreamer = (*Branch)(nil)

// EOF
