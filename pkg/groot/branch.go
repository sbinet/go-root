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
	
	basketBytes []int32 // length of baskets on file
	basketEntry []int32 // table of first entry of each basket
	basketSeek  []int64 // addresses of baskets on file
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
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	dprintf("vers=%v spos=%v pos=%v bcnt=%v\n", vers, spos, pos, bcnt)
	branch.name, branch.title = b.read_tnamed()
	dprintf("name='%v' title='%v'\n", branch.name, branch.title)
	dprintf("spos=%v\n", b.Pos())

	maxbaskets := uint32(0)
	splitlvl := int32(0)
	if vers <= 5 {
		b.ntoi4() // fCompress
		b.ntoi4() // fBasketSize
		branch.entryOffsetLen = b.ntou4()
		maxbaskets = b.ntou4() // fMaxBaskets
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
		maxbaskets = b.ntou4() // fMaxBaskets
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
		maxbaskets = b.ntou4() // fMaxBaskets
		splitlvl = b.ntoi4() // fSplitLevel
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
		maxbaskets = b.ntou4() // fMaxBaskets
		splitlvl = b.ntoi4() // fSplitLevel
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
		maxbaskets = b.ntou4() // fMaxBaskets
		splitlvl = b.ntoi4() // fSplitLevel
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
		maxbaskets = b.ntou4() // fMaxBaskets
		splitlvl = b.ntoi4() // fSplitLevel
		b.ntou8() // entries
		b.ntou8() // fFirstEntry
		b.ntou8() // tot_bytes
		b.ntou8() // zip_bytes
	}
	dprintf("::branch::stream : [%s] split-lvl= %v\n", branch.name, splitlvl)

	dprintf("::branch::stream : branches : begin\n")
	branches := b.read_obj_array()
	dprintf("::branch::stream : branches : end\n")
	dprintf("sub-branches: %v\n", len(branches))

	dprintf("::branch::stream : leaves : begin\n")
	leaves := b.read_obj_array()
	dprintf("::branch::stream : leaves : end\n")
	dprintf("sub-leaves: %v\n", len(leaves))

	dprintf("::branch::stream : streamed_baskets : begin\n")
	baskets := b.read_obj_array()
	dprintf("::branch::stream : streamed_baskets : end\n")
	dprintf("baskets: %v\n", len(baskets))

	branch.basketEntry = make([]int32, 0, int(maxbaskets))
	branch.basketBytes = make([]int32, 0, int(maxbaskets))
	branch.basketSeek =  make([]int64, int(maxbaskets))

	if vers < 6 {
		copy(branch.basketEntry, b.read_array_I())
		if vers <= 4 {
			branch.basketBytes = make([]int32, int(maxbaskets))
		} else {
			copy(branch.basketBytes, b.read_array_I())
		}
		if vers < 2 {
			panic("branch.ROOTDecode of vers<2 *NOT* handled")
		} else {
			nseeks := int(b.ntoi4())
			for i := 0; i<nseeks; i++ {
				branch.basketSeek[i] = int64(b.ntoi4())
			}
		}
	} else if vers <= 9 {
		// see TStreamerInfo::ReadBuffer::ReadBasicPointer
		isarray := byte(0)
		isarray = b.ntobyte()
		if isarray != 0 {
			copy(branch.basketBytes, b.read_fast_array_I(int(maxbaskets)))
		}
		isarray = b.ntobyte()
		if isarray != 0 {
			copy(branch.basketEntry, b.read_fast_array_I(int(maxbaskets)))
		}
		isbigfile := b.ntobyte()
		if isbigfile == 2 {
			copy(branch.basketSeek, b.read_fast_array_L(int(maxbaskets)))
		} else {
			for i := 0; i<int(maxbaskets); i++ {
				branch.basketSeek[i] = int64(b.ntoi4())
			}
			
		}
	} else { // vers >= 10
		// see TStreamerInfo::ReadBuffer::ReadBasicPointer
		isarray := byte(0)

		isarray = b.ntobyte()
		if isarray != 0 {
			copy(branch.basketBytes, b.read_fast_array_I(int(maxbaskets)))
		}

		isarray = b.ntobyte()
		if isarray != 0 {
			bentries := b.read_fast_array_UL(int(maxbaskets))
			for _,v := range bentries {
				branch.basketEntry = append(branch.basketEntry, int32(v))
			}
		}

		isarray = b.ntobyte()
		if isarray != 0 {
			bentries := b.read_fast_array_UL(int(maxbaskets))
			for i,v := range bentries {
				branch.basketSeek[i] = int64(v)
			}
		}
	}

	if vers > 2 {
		fname := b.read_tstring()
		dprintf("fname=%s\n", fname)
	}
	dprintf("spos=%v\n", b.Pos())
	b.check_byte_count(pos, bcnt, spos, "TBranch")
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
