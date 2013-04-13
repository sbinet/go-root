package groot

import (
	"fmt"
	"reflect"
)

type Tree struct {
	file      *File
	name      string
	title     string
	entries   uint64
	tot_bytes uint64
	zip_bytes uint64
	branches  []Branch
}

func NewTree(file *File, name, title string) (tree *Tree, err error) {
	tree = &Tree{
		file:     file,
		name:     name,
		title:    title,
		entries:  0,
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

func (tree *Tree) Class() string {
	return "TTree"
}

func (tree *Tree) Name() string {
	return tree.name
}

func (tree *Tree) Title() string {
	return tree.title
}

func (tree *Tree) Entries() uint64 {
	return tree.entries
}

func (tree *Tree) Branches() []Branch {
	return tree.branches
}

func (tree *Tree) ROOTDecode(b *Buffer) (err error) {
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	printf("[tree] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	tree.name, tree.title = b.read_tnamed()
	printf("name='%v' title='%v'\n", tree.name, tree.title)
	b.read_attline()
	b.read_attfill()
	b.read_attmarker()

	if vers <= 4 {
		b.ntoi4() //fScanField
		b.ntoi4() //fMaxEntryLoop
		b.ntoi4() //fMaxVirtualSize
		tree.entries = uint64(b.ntod())
		tree.tot_bytes = uint64(b.ntod())
		tree.zip_bytes = uint64(b.ntod())
		b.ntoi4() //fAutoSave
		b.ntoi4() //fEstimate
	} else if vers <= 9 {
		tree.entries = uint64(b.ntod())
		tree.tot_bytes = uint64(b.ntod())
		tree.zip_bytes = uint64(b.ntod())
		b.ntod()  //fSaveBytes
		b.ntoi4() //fTimerInterval
		b.ntoi4() //fScanField
		b.ntoi4() //fUpdate
		b.ntoi4() //fMaxEntryLoop
		b.ntoi4() //fMaxVirtualSize
		b.ntoi4() //fAutoSave
		b.ntoi4() //fEstimate
	} else if vers < 16 { //FIXME: what is the exact version ?
		tree.entries = uint64(b.ntod())
		tree.tot_bytes = uint64(b.ntod())
		tree.zip_bytes = uint64(b.ntod())
		b.ntod()  //fSaveBytes
		b.ntod()  //fWeight
		b.ntoi4() //fTimerInterval
		b.ntoi4() //fScanField
		b.ntoi4() //fUpdate
		b.ntoi4() //fMaxEntryLoop
		b.ntoi4() //fMaxVirtualSize
		b.ntoi4() //fAutoSave
		b.ntoi4() //fEstimate
	} else { // vers >= 16
		tree.entries = b.ntou8()
		tree.tot_bytes = b.ntou8()
		tree.zip_bytes = b.ntou8()
		b.ntou8() //fSavedBytes
		if vers >= 18 {
			b.ntoi8() //fFlushedBytes
		}
		b.ntod()  //fWeight
		b.ntoi4() //fTimerInterval
		b.ntoi4() //fScanField
		b.ntoi4() //fUpdate
		if vers >= 18 {
			b.ntoi4() //fDefaultEntryOffsetLen
		}
		b.ntoi8() //fMaxEntries
		b.ntoi8() //fMaxEntryLoop
		b.ntou8() //fMaxVirtualSize
		b.ntou8() //fAutoSave
		if vers >= 18 {
			b.ntoi8() //fAutoFlush
		}
		b.ntoi8() //fEstimate
	}

	printf("=> (%s) entries=%v tot_bytes=%v zip_bytes=%v\n",
		tree.name, tree.entries, tree.tot_bytes, tree.zip_bytes)

	branches := b.read_obj_array()
	printf("-- #nbranches: %v\n", len(branches))
	tree.branches = make([]Branch, len(branches))
	for i, v := range branches {
		tree.branches[i] = *(v.(ibranch).toBranch())
	}
	leaves := b.read_obj_array()
	printf("-- #nleaves: %v\n", len(leaves))

	if vers >= 10 {
		b.read_object() // fAliases *TList
	}

	b.read_array_D() // fIndexValues TArrayD
	b.read_array_I() // fIndex TArrayI

	if vers >= 16 {
		b.read_object() // fTreeIndex *TVirtualIndex // FIXME?
	}
	if vers >= 6 {
		b.read_object() // fFriends *TList
	}
	if vers >= 16 {
		b.read_object() // fUserInfo *TList
		b.read_object() // fBranchRef *TBranchRef
	}

	b.check_byte_count(pos, bcnt, spos, "TTree")
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
