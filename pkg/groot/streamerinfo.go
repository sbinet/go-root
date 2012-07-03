package groot

import (
	"reflect"
)

type StreamerInfo struct {
	name      string
	title     string
	checksum  uint32
	classvers uint32
	elmts     []Object
}

func (si *StreamerInfo) Class() string {
	return "TStreamerInfo"
}

func (si *StreamerInfo) Name() string {
	return si.name
}

func (si *StreamerInfo) Title() string {
	return si.title
}

func (si *StreamerInfo) ROOTDecode(b *Buffer) (err error) {

	vers, pos, bcnt := b.read_version()
	dprintf("vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	si.name, si.title = b.read_tnamed()
	dprintf("name='%v' title='%v'\n", si.name, si.title)
	si.checksum = b.ntou4()
	si.classvers = b.ntou4()
	si.elmts = b.read_elements()

	return
}

func (si *StreamerInfo) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

func init() {
	f := func() reflect.Value {
		o := &StreamerInfo{}
		return reflect.ValueOf(o)
	}
	Factory.db["TStreamerInfo"] = f
	Factory.db["*groot.StreamerInfo"] = f
}

// check interfaces
var _ Object = (*StreamerInfo)(nil)
var _ ROOTStreamer = (*StreamerInfo)(nil)

// EOF
