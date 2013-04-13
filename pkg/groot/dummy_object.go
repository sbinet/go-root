package groot

import (
	"reflect"
	"unsafe"
)

type dummyObject struct {
}

func (d *dummyObject) Class() string {
	return "dummy-object-type"
}

func (d *dummyObject) Name() string {
	return "dummy-object-name"
}

func (d *dummyObject) Title() string {
	return "dummy-object-title"
}

func (d *dummyObject) ROOTDecode(b *Buffer) (err error) {

	spos := b.Pos()
	vers, pos, bcnt := b.clone().read_version()
	printf("dummy: vers=%v spos=%v pos=%v bcnt=%v\n", vers, spos, pos, bcnt)
	b.read_nbytes(int(bcnt) + int(unsafe.Sizeof(uint(0))))
	b.check_byte_count(pos, bcnt, spos, "dummy")
	printf("dummy: vers=%v spos=%v pos=%v bcnt=%v [done]\n", vers, spos, pos, bcnt)
	return
}

func (d *dummyObject) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

func init() {
	f := func() reflect.Value {
		o := &dummyObject{}
		return reflect.ValueOf(o)
	}
	Factory.db["*groot.dummyObject"] = f
}

// check interfaces
var _ Object = (*dummyObject)(nil)
var _ ROOTStreamer = (*dummyObject)(nil)
