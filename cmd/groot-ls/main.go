// groot-ls recursively dumps the hierarchy tree of a ROOT file
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/binet/go-root/pkg/groot"
)

var fname = flag.String("f", "ntuple.root", "ROOT file to inspect")

func normpath(path []string) string {
	name := strings.Join(path, "/")
	if len(name) > 2 && name[:2] == "//" {
		name = name[1:]
	}
	return name
}

func main() {
	fmt.Printf(":: groot-ls ::\n")
	flag.Parse()

	nfactories := groot.Factory.NumKey()
	fmt.Printf(":: groot has registered %d factory-type(s)\n", nfactories)

	f, err := groot.NewFileReader(*fname)
	if err != nil {
		fmt.Printf("**error**: %v\n", err)
		os.Exit(1)
	}

	if f == nil {
		fmt.Printf("**error**: invalid file pointer\n")
		os.Exit(1)
	}

	fmt.Printf("file: '%s' (version=%v)\n", f.Name(), f.Version())

	var inspect func(*groot.Directory, []string, string)
	inspect = func(dir *groot.Directory, path []string, indent string) {
		name := normpath(path)
		if dir == nil {
			fmt.Printf("err: invalid directory [%s]\n", name)
			return
		}
		keys := dir.Keys()
		fmt.Printf("%s%s -> #%d key(s)\n", indent, name, len(keys))
		for _, k := range keys {
			fmt.Printf("%skey: name='%s' title='%s' type=%s\n",
				indent, k.Name(), k.Title(), k.Class())
			if v, ok := k.Value().(*groot.Directory); ok {
				path := append(path, k.Name())
				inspect(v, path, indent+"  ")
			}
		}
	}
		
	dir := f.Dir()
	inspect(dir, []string{"/"}, "")

	fmt.Printf("::bye.\n")
}

// EOF
