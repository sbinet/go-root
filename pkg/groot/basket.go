package groot

import (
	"fmt"
	"reflect"
)

type Basket struct {
	key Key

	nev_bufsz    uint32 // Length in Int_t of entry_offset
	nev          uint32 // Number of entries in basket
	last         uint32 // Pointer to last used byte in basket
	entry_offset int    // [m_nev] Offset of entries in Key.buffer
	displacement int    //![m_nev] Displacement of entries in Key.buffer

}

func (basket *Basket) Class() Class {
	//FIXME
	panic("not implemented")
}

func (basket *Basket) Name() string {
	return basket.key.Name()
}

func (basket *Basket) Title() string {
	return basket.key.Title()
}

func (basket *Basket) ROOTDecode(b *Buffer) (err error) {
	startpos := b.Len()
	k,err := NewKey(nil, 0, 0)
	if err != nil {
		return err
	}
	basket.key = *k
	err = basket.key.init_from_buffer(b)
	if err != nil {
		return err
	}

	vers, pos, bcnt := b.read_version()
	dprintf("vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)
	bufsz := b.ntou4()
	basket.nev_bufsz = b.ntou4()
	basket.nev = b.ntou4()
	basket.last = b.ntou4()
	flag := b.ntobyte()
	if basket.last > bufsz {
		bufsz = basket.last
	}
	
	basket_key_len := b.Len() - startpos
	if basket_key_len != int(basket.key.keysz) {
		basket.key.keysz = uint16(basket_key_len)
	}
	if basket.key.objsz != 0 {
		basket.key.objsz = uint32(basket.key.nbytes) - uint32(basket.key.keysz)
	}
	if flag == 0 {
		// fHeaderOnly
		return
	}
	// from Guy: adding this useful (?) test
	if (flag!=1 && flag!=2 &&
		flag!=11 && flag!=12 &&
		flag!=41 && flag!=42 &&
		flag!=51 && flag!=52) {
		err = fmt.Errorf("groot.basket.ROOTDecode: bad flag (=%v)",
			int(flag))
	}
	return
}

func (basket *Basket) ROOTEncode(b *Buffer) (err error) {
	//bb := b.Clone()
	//err = basket.key.ROOTEncode(b)
	if err != nil {
		return err
	}
	panic("fixme")
	return
}

func init() {
	f := func() reflect.Value {
		o := &Basket{}
		return reflect.ValueOf(o)
	}
	Factory.db["TBasket"] = f
	Factory.db["*groot.Basket"] = f
}

// EOF
