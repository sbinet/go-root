package groot

import (
	"reflect"
)

type Basket struct {
	Key

	nev_bufsz    uint32 // Length in Int_t of entry_offset
	nev          uint32 // Number of entries in basket
	last         uint32 // Pointer to last used byte in basket
	entry_offset int    // [m_nev] Offset of entries in Key.buffer
	displacement int    //![m_nev] Displacement of entries in Key.buffer

}

/*
func (basket *Basket) ROOTDecode(b *Buffer) (err error) {
	err = basket.Key.ROOTDecode(b)
	if err != nil {
		return err
	}
	panic("fixme")
	return
}

func (basket *Basket) ROOTEncode(b *Buffer) (err error) {
	err = basket.Key.ROOTEncode(b)
	if err != nil {
		return err
	}
	panic("fixme")
	return
}
*/

func init() {
	f := func() reflect.Value {
		o := &Basket{}
		return reflect.ValueOf(o)
	}
	Factory.db["TBasket"] = f
	Factory.db["*groot.Basket"] = f
}

// EOF
