package groot

type ibaseLeaf interface {
	toBaseLeaf() *baseLeaf
}

type baseLeaf struct {
	name   string
	title  string
	ndata  uint32 // number of elements in fAddress
	length uint32 // number of fixed length elements

	leaf_count *baseLeaf // pointer to Leaf-count if variable length
}

func (base *baseLeaf) Class() Class {
	panic("not implemented")
}

func (base *baseLeaf) Name() string {
	return base.name
}

func (base *baseLeaf) Title() string {
	return base.title
}

func (base *baseLeaf) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()

	vers, pos, bcnt := b.read_version()
	printf("baseleaf-vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	base.name, base.title = b.read_tnamed()
	printf("baseleaf-name='%v' title='%v'\n", base.name, base.title)
	base.length = b.ntou4()
	printf("baseleaf-length=%v\n", base.length)
	b.ntoi4()   // fLengthType
	b.ntoi4()   // fOffset
	b.ntobyte() // fIsRange
	b.ntobyte() // fIsUnsigned

	obj := b.read_object()
	printf("baseleaf-nobjs: %v\n", obj)
	if obj != nil {
		base.leaf_count = obj.(ibaseLeaf).toBaseLeaf()
	}

	if base.length == 0 {
		//FIXME: ??? really ??? (check with Guy)
		base.length = 1
	}

	b.check_byte_count(pos, bcnt, spos, "TLeaf")
	return
}

// func init() {
// 	f := func() reflect.Value {
// 		o := &BaseLeaf{}
// 		return reflect.ValueOf(o)
// 	}
// 	Factory.db["TBaseLeaf"] = f
// 	Factory.db["*groot.BaseLeaf"] = f
// }

// EOF
