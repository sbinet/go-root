package groot

import (
	"reflect"
)

type BranchElement struct {
	branch Branch
	object Object
	class  string // class name of referenced object
	vers   int    // version number of class
	id     int    // element serial number in fInfo
	btype  int    // branch type
	stype  int    // branch streamer type
}

func (be *BranchElement) toBranch() *Branch {
	return &be.branch
}

func (be *BranchElement) Class() string {
	panic("not implemented")
}

func (be *BranchElement) Name() string {
	return be.branch.name
}

func (be *BranchElement) Title() string {
	return be.branch.title
}

func (be *BranchElement) ROOTDecode(b *Buffer) (err error) {

	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("[branch_element] vers=%v spos=%v pos=%v bcnt=%v\n",
		vers, spos, pos, bcnt)
	err = be.branch.ROOTDecode(b)
	if err != nil {
		return err
	}

	if vers <= 7 {
		be.class = b.read_tstring()
		be.vers = int(b.ntoi4())
		be.id = int(b.ntoi4())
		be.btype = int(b.ntoi4())
		be.stype = int(b.ntoi4())
	} else { // vers >= 8
		be.class = b.read_tstring()
		b.read_tstring() // fParentName
		b.read_tstring() // fCloneName
		b.ntoi4()        // fCheckSum
		be.vers = int(b.ntoi4())
		be.id = int(b.ntoi4())
		be.btype = int(b.ntoi4())
		be.stype = int(b.ntoi4())
		b.ntoi4() // fMaximum

		b.read_object() // fBranchCount  *TBranchElement
		b.read_object() // fBranchCount2 *TBranchElement
	}
	b.check_byte_count(pos, bcnt, spos, "TBranchElement")
	return
}

func (be *BranchElement) ROOTEncode(b *Buffer) (err error) {
	panic("not implemented")
	return
}

func init() {
	f := func() reflect.Value {
		o := &BranchElement{}
		return reflect.ValueOf(o)
	}
	Factory.db["TBranchElement"] = f
	Factory.db["*groot.BranchElement"] = f
}

// check interfaces
var _ Object = (*BranchElement)(nil)
var _ ROOTStreamer = (*BranchElement)(nil)
