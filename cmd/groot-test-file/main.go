package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/binet/go-root/pkg/groot"
)

var fname = flag.String("f", "small.ntuple.0.root", "ROOT file to inspect")

func main() {
	fmt.Printf("== test go-root ==\n")
	flag.Parse()

	f, err := groot.NewFileReader(*fname)
	if err != nil {
		fmt.Printf("**error**: %v\n", err)
		os.Exit(1)
	}

	if f == nil {
		fmt.Printf("invalid file pointer\n")
		os.Exit(1)
	}

	fmt.Printf("f: %s (version=%v)\n", f.Name(), f.Version())

	var fct func(*groot.Directory, string, string)
	fct = func(dir *groot.Directory, name, indent string) {
		if dir == nil {
			fmt.Printf("err: invalid directory [%s]\n", name)
			return
		}
		keys := dir.Keys()
		fmt.Printf("%s[%s] -> #%d key(s)\n", indent, name, len(keys))
		for i, k := range keys {
			fmt.Printf("%skey[%d]: [name=%s] [title=%s] [type=%s]\n",
				indent, i, k.Name(), k.Title(), k.Class())
			bufkey, err := k.Buffer()
			fmt.Printf("buf: %d, (err=%v)\n", len(bufkey), err)
			if k.Class() == "TDirectory" {
				buf, err := k.TBuffer()
				if err != nil {
					fmt.Printf("**err**: %v\n", err)
					return
				}
				v, err := groot.NewDirectory(f, buf.Buffer())
				if err != nil {
					fmt.Printf("**err**: %v\n", err)
					return
				}
				fct(v, name+"/"+k.Name(), indent+"  ")
			}
		}
	}
		
	dir := f.Dir()
	fct(dir, "/", "")

	fmt.Printf("::bye.\n")
}

// EOF
