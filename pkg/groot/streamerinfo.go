package groot

import (
	"fmt"
	"reflect"
)

type StreamerInfo struct {
	name      string
	title     string
	checksum  uint32
	classvers uint32
	elmts     []StreamerElement
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

	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerinfo] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	si.name, si.title = b.read_tnamed()
	printf("name='%v' title='%v'\n", si.name, si.title)
	if vers <= 1 {
		err = fmt.Errorf("too old version for StreamerInfo (v=%d)", vers)
		return
	}
	si.checksum = b.ntou4()
	si.classvers = b.ntou4()
	elmts := b.read_elements()
	si.elmts = make([]StreamerElement, 0, len(elmts))
	for _, v := range elmts {
		switch vv := v.(type) {
		case StreamerElement:
			si.elmts = append(si.elmts, vv)
		default:
			si.elmts = append(si.elmts, nil)
		}
	}

	b.check_byte_count(pos, bcnt, spos, "TStreamerInfo")
	return
}

func (si *StreamerInfo) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

type StreamerElement interface {
	Name() string
	Title() string
	Type() int       // element type
	Size() int       // sizeof element
	ArrLen() int     // cumulative size of all array dims
	ArrDim() int     // number of array dimensions
	MaxIdx() []int32 // maximum array index for array dimension "dim"
	Offset() int     // element offset in class
	//IsNewType() int // new element type when reading
	TypeName() string // data type name of data member
}

type seBase struct {
	name  string
	title string

	etype    int     // element type
	esize    int     // sizeof element
	arrlen   int     // cumulative size of all array dims
	arrdim   int     // number of array dimensions
	maxidx   []int32 // maximum array index for array dimension "dim"
	offset   int     // element offset in class
	newtype  int     // new element type when reading
	typename string  // data type name of data member
}

func (se *seBase) Class() string {
	return "TStreamerElement"
}

func (se *seBase) Name() string {
	return se.name
}

func (se *seBase) Title() string {
	return se.title
}

func (se *seBase) Type() int {
	return se.etype
}

func (se *seBase) Size() int {
	return se.esize
}

func (se *seBase) ArrLen() int {
	return se.arrlen
}

func (se *seBase) ArrDim() int {
	return se.arrdim
}

func (se *seBase) MaxIdx() []int32 {
	return se.maxidx
}

func (se *seBase) Offset() int {
	return se.offset
}

func (se *seBase) TypeName() string {
	return se.typename
}

func (se *seBase) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerelmt] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerElement.ROOTDecode: version<2 is not supported")
	}
	se.name, se.title = b.read_tnamed()
	se.etype = int(b.ntoi4())
	se.esize = int(b.ntoi4())
	se.arrlen = int(b.ntoi4())
	se.arrdim = int(b.ntoi4())
	se.maxidx = b.read_fast_array_I(5) // FIXME: magic constant
	se.typename = b.read_tstring()

	b.check_byte_count(pos, bcnt, spos, "TStreamerElement")
	return
}

func (se *seBase) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerBase is a streamer element for a base class
type StreamerBase struct {
	seBase
	version int // version number of the base class
}

func (se *StreamerBase) Class() string {
	return "TStreamerBase"
}

func (se *StreamerBase) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerbase] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerBase.ROOTDecode: version<2 is not supported")
	}
	err = se.seBase.ROOTDecode(b)
	if err != nil {
		return err
	}
	se.version = int(b.ntoi4())
	b.check_byte_count(pos, bcnt, spos, "TStreamerBase")
	return
}

func (se *StreamerBase) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerBasicType is a streamer element for a builtin type
type StreamerBasicType struct {
	seBase
}

func (se *StreamerBasicType) Class() string {
	return "TStreamerBasicType"
}

func (se *StreamerBasicType) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerbasictype] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerBasicType.ROOTDecode: version<2 is not supported")
	}
	err = se.seBase.ROOTDecode(b)
	if err != nil {
		return err
	}
	b.check_byte_count(pos, bcnt, spos, "TStreamerBasicType")
	return
}

func (se *StreamerBasicType) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerBasicPointer is a streamer element for a pointer to a builtin type
type StreamerBasicPointer struct {
	seBase
	countvers  int    // version number of the class with the counter
	countname  string // name of the data member holding the array count
	countclass string // name of the class with the counter
}

func (se *StreamerBasicPointer) Class() string {
	return "TStreamerBasicPointer"
}

func (se *StreamerBasicPointer) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerbasicptr] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerBasicPointer.ROOTDecode: version<2 is not supported")
	}
	err = se.seBase.ROOTDecode(b)
	if err != nil {
		return err
	}

	se.countvers = int(b.ntoi4())
	se.countname = b.read_tstring()
	se.countclass = b.read_tstring()
	printf("[streamerbasicptr] cntvers=%v name=%v cls=%v\n",
		se.countvers, se.countname, se.countclass)

	b.check_byte_count(pos, bcnt, spos, "TStreamerBasicPointer")
	return
}

func (se *StreamerBasicPointer) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerString is a streamer element for a string
type StreamerString struct {
	seBase
}

func (se *StreamerString) Class() string {
	return "TStreamerString"
}

func (se *StreamerString) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerstring] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerString.ROOTDecode: version<2 is not supported")
	}
	err = se.seBase.ROOTDecode(b)
	if err != nil {
		return err
	}
	b.check_byte_count(pos, bcnt, spos, "TStreamerString")
	return
}

func (se *StreamerString) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerObject is a streamer element for an object
type StreamerObject struct {
	seBase
}

func (se *StreamerObject) Class() string {
	return "TStreamerObject"
}

func NewStreamerObject(name, title string, offset int, typename string) *StreamerObject {
	o := &StreamerObject{
		seBase: seBase{
			name:     name,
			title:    title,
			offset:   offset,
			etype:    kObject,
			typename: typename,
		},
	}
	switch name {
	case "TObject":
		o.seBase.etype = kTObject
	case "TNamed":
		o.seBase.etype = kTNamed
	}
	return o
}

func (se *StreamerObject) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerobj] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerObject.ROOTDecode: version<2 is not supported")
	}
	err = se.seBase.ROOTDecode(b)
	if err != nil {
		return err
	}
	b.check_byte_count(pos, bcnt, spos, "TStreamerObject")
	return
}

func (se *StreamerObject) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerObjectPointer is a streamer element for a pointer to an object
type StreamerObjectPointer struct {
	seBase
}

func (se *StreamerObjectPointer) Class() string {
	return "TStreamerObjectPointer"
}

func (se *StreamerObjectPointer) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerobjptr] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerObjectPointer.ROOTDecode: version<2 is not supported")
	}
	err = se.seBase.ROOTDecode(b)
	if err != nil {
		return err
	}
	b.check_byte_count(pos, bcnt, spos, "TStreamerObjectPointer")
	return
}

func (se *StreamerObjectPointer) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerObjectAny is a streamer element for any object
type StreamerObjectAny struct {
	seBase
}

func (se *StreamerObjectAny) Class() string {
	return "TStreamerObjectAny"
}

func (se *StreamerObjectAny) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerobjany] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerObjectAny.ROOTDecode: version<2 is not supported")
	}
	err = se.seBase.ROOTDecode(b)
	if err != nil {
		return err
	}
	b.check_byte_count(pos, bcnt, spos, "TStreamerObjectAny")
	return
}

func (se *StreamerObjectAny) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerSTL is a streamer element for STL containers
type StreamerSTL struct {
	seBase
	stltype int // type of STL container
	ctype   int // type of contained object
}

func (se *StreamerSTL) Class() string {
	return "TStreamerSTL"
}

func (se *StreamerSTL) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerstl] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerSTL.ROOTDecode: version<2 is not supported")
	}
	err = se.seBase.ROOTDecode(b)
	if err != nil {
		return err
	}
	se.stltype = int(b.ntoi4())
	se.ctype = int(b.ntoi4())
	printf("[streamerstl] name='%v' title='%s' type=%d stl=%v ctype=%v\n",
		se.Name(), se.Title(), se.Type(), se.stltype, se.ctype)
	b.check_byte_count(pos, bcnt, spos, "TStreamerSTL")
	return
}

func (se *StreamerSTL) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// StreamerSTLstring is a streamer element for std::string
type StreamerSTLstring struct {
	StreamerSTL
}

func (se *StreamerSTLstring) Class() string {
	return "TStreamerSTLstring"
}

func (se *StreamerSTLstring) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("[streamerstlstr] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	if vers < 2 {
		panic("groot.StreamerSTLstring.ROOTDecode: version<2 is not supported")
	}
	err = se.StreamerSTL.ROOTDecode(b)
	if err != nil {
		return err
	}
	printf("[streamerstlstr] name='%v' title='%s' type=%d stl=%v ctype=%v\n",
		se.Name(), se.Title(), se.Type(), se.stltype, se.ctype)
	b.check_byte_count(pos, bcnt, spos, "TStreamerSTLstring")
	return
}

func (se *StreamerSTLstring) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

// register factories ---------------------------------------------------------

func init() {

	{
		f := func() reflect.Value {
			o := &StreamerInfo{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerInfo"] = f
		Factory.db["*groot.StreamerInfo"] = f
	}

	{
		f := func() reflect.Value {
			o := &seBase{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerElement"] = f
		Factory.db["groot.StreamerElement"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerBase{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerBase"] = f
		Factory.db["*groot.StreamerBase"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerBasicType{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerBasicType"] = f
		Factory.db["*groot.StreamerBasicType"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerBasicPointer{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerBasicPointer"] = f
		Factory.db["*groot.StreamerBasicPointer"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerString{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerString"] = f
		Factory.db["*groot.StreamerString"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerObject{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerObject"] = f
		Factory.db["*groot.StreamerObject"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerObjectPointer{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerObjectPointer"] = f
		Factory.db["*groot.StreamerObjectPointer"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerObjectAny{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerObjectAny"] = f
		Factory.db["*groot.StreamerObjectAny"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerSTL{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerSTL"] = f
		Factory.db["*groot.StreamerSTL"] = f
	}

	{
		f := func() reflect.Value {
			o := &StreamerSTLstring{}
			return reflect.ValueOf(o)
		}
		Factory.db["TStreamerSTLstring"] = f
		Factory.db["*groot.StreamerSTLstring"] = f
	}
}

// check interfaces -----------------------------------------------------------
var _ Object = (*StreamerInfo)(nil)
var _ ROOTStreamer = (*StreamerInfo)(nil)

var _ Object = (*seBase)(nil)
var _ ROOTStreamer = (*seBase)(nil)
var _ StreamerElement = (*seBase)(nil)

var _ Object = (*StreamerBase)(nil)
var _ ROOTStreamer = (*StreamerBase)(nil)
var _ StreamerElement = (*StreamerBase)(nil)

var _ Object = (*StreamerBasicType)(nil)
var _ ROOTStreamer = (*StreamerBasicType)(nil)
var _ StreamerElement = (*StreamerBasicType)(nil)

var _ Object = (*StreamerBasicPointer)(nil)
var _ ROOTStreamer = (*StreamerBasicPointer)(nil)
var _ StreamerElement = (*StreamerBasicPointer)(nil)

var _ Object = (*StreamerString)(nil)
var _ ROOTStreamer = (*StreamerString)(nil)
var _ StreamerElement = (*StreamerString)(nil)

var _ Object = (*StreamerObject)(nil)
var _ ROOTStreamer = (*StreamerObject)(nil)
var _ StreamerElement = (*StreamerObject)(nil)

var _ Object = (*StreamerObjectPointer)(nil)
var _ ROOTStreamer = (*StreamerObjectPointer)(nil)
var _ StreamerElement = (*StreamerObjectPointer)(nil)

var _ Object = (*StreamerObjectAny)(nil)
var _ ROOTStreamer = (*StreamerObjectAny)(nil)
var _ StreamerElement = (*StreamerObjectAny)(nil)

var _ Object = (*StreamerSTL)(nil)
var _ ROOTStreamer = (*StreamerSTL)(nil)
var _ StreamerElement = (*StreamerSTL)(nil)

var _ Object = (*StreamerSTLstring)(nil)
var _ ROOTStreamer = (*StreamerSTLstring)(nil)
var _ StreamerElement = (*StreamerSTLstring)(nil)

// EOF
