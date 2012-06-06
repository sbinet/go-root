package groot

import (
	"reflect"
)

type FactoryFct func() reflect.Value

type factory struct {
	db map[string]FactoryFct // a registry of all factory functions by class name
}

func (f *factory) NumKey() int {
	return len(f.db)
}

func (f *factory) Keys() []string {
	keys := make([]string, 0, len(f.db))
	for k, _ := range f.db {
		keys = append(keys, k)
	}
	return keys
}

func (f *factory) HasKey(n string) bool {
	_, ok := f.db[n]
	return ok
}

func (f *factory) Get(n string) FactoryFct {
	fct, ok := f.db[n]
	if ok {
		return fct
	}
	return nil
}

// the registry of all factory functions, by class name
var Factory factory

func init() {
	Factory = factory{
		db: make(map[string]FactoryFct),
	}

	Factory.db["TDirectory"] = func() reflect.Value {
		o := &Directory{file: nil, keys: make([]Key, 0)}
		return reflect.ValueOf(o)
	}
}

// EOF
