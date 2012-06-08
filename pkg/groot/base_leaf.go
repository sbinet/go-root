package groot

import (
	"reflect"
)

type BaseLeaf struct {
	name   string
	title  string
	ndata  uint32 // number of elements in fAddress
	length uint32 // number of fixed length elements

	leaf_count *BaseLeaf // pointer to Leaf-count if variable length
}

func init() {
	f := func() reflect.Value {
		o := &BaseLeaf{}
		return reflect.ValueOf(o)
	}
	Factory.db["TBaseLeaf"] = f
	Factory.db["*groot.BaseLeaf"] = f
}

// EOF
