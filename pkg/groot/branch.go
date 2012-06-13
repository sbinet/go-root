package groot

import (
	"reflect"
)

type Branch struct {
	name  string
	title string

	file *File

	autodelete bool
	branches []Branch
	leaves   []baseLeaf
	baskets  []Basket
	entryOffsetLen uint32 // initial length of fEntryOffset table in the basket buffers
	writeBasket uint32 // last basket number written
	entryNumber uint32 // current entry number (last one filled in this branch)
	readBasket uint32  // current basket number when reading
	
	//fBasketBytes
	//fBasketEntry
	//fBasketSeek

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

	if vers <= 5 {
		b.ntoi4() // fCompress
		b.ntoi4() // fBasketSize
		branch.entryOffsetLen = b.ntou4()
		b.ntou4() // fMaxBaskets
		branch.writeBasket = b.ntou4()
		branch.entryNumber = b.ntou4()
		b.ntod() // entries
		b.ntod() // tot_bytes
		b.ntod() // zip_bytes
		b.ntoi4() // fOffset
	} else if vers <= 6 {
		b.ntoi4() // fCompress
		b.ntoi4() // fBasketSize
		branch.entryOffsetLen = b.ntou4()
		branch.writeBasket = b.ntou4()
		branch.entryNumber = b.ntou4()
		b.ntoi4() // fOffset
		b.ntou4() // fMaxBaskets
		b.ntod() // entries
		b.ntod() // tot_bytes
		b.ntod() // zip_bytes
	} else if vers <= 7 {
		b.ntoi4() // fCompress
		b.ntoi4() // fBasketSize
		branch.entryOffsetLen = b.ntou4()
		branch.writeBasket = b.ntou4()
		branch.entryNumber = b.ntou4()
		b.ntoi4() // fOffset
		b.ntou4() // fMaxBaskets
		b.ntoi4() // fSplitLevel
		b.ntod() // entries
		b.ntod() // tot_bytes
		b.ntod() // zip_bytes
	} else if vers <= 9 {
		b.read_attfill()
		b.ntoi4() // fCompress
		b.ntoi4() // fBasketSize
		branch.entryOffsetLen = b.ntou4()
		branch.writeBasket = b.ntou4()
		branch.entryNumber = b.ntou4()
		b.ntoi4() // fOffset
		b.ntou4() // fMaxBaskets
		b.ntoi4() // fSplitLevel
		b.ntod() // entries
		b.ntod() // tot_bytes
		b.ntod() // zip_bytes
	} else if vers <= 10 {
		b.read_attfill()
		b.ntoi4() // fCompress
		b.ntoi4() // fBasketSize
		branch.entryOffsetLen = b.ntou4()
		branch.writeBasket = b.ntou4()
		branch.entryNumber = uint32(b.ntou8()) //fixme ?
		b.ntoi4() // fOffset
		b.ntou4() // fMaxBaskets
		b.ntoi4() // fSplitLevel
		b.ntou8() // entries
		b.ntou8() // tot_bytes
		b.ntou8() // zip_bytes
	} else { //vers>=11
		b.read_attfill()
		b.ntoi4() // fCompress
		b.ntoi4() // fBasketSize
		branch.entryOffsetLen = b.ntou4()
		branch.writeBasket = b.ntou4()
		branch.entryNumber = uint32(b.ntou8()) //fixme ?
		b.ntoi4() // fOffset
		b.ntou4() // fMaxBaskets
		b.ntoi4() // fSplitLevel
		b.ntou8() // entries
		b.ntou8() // fFirstEntry
		b.ntou8() // tot_bytes
		b.ntou8() // zip_bytes
	}
	branches := b.read_obj_array()
	dprintf("sub-branches: %v\n", len(branches))

	leaves := b.read_obj_array()
	dprintf("sub-leaves: %v\n", len(leaves))

	baskets := b.read_obj_array()
	dprintf("baskets: %v\n", len(baskets))

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
