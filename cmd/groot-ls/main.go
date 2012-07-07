// groot-ls recursively dumps the hierarchy tree of a ROOT file
package main

import (
	"flag"
	"fmt"
	//"log"
	"os"
	//"runtime/pprof"
	"strings"

	"github.com/sbinet/go-root/pkg/groot"
)

var fname = flag.String("f", "", "ROOT file to inspect")
var detailed = flag.Bool("detailed", false, "enable detailed dump (of trees)")
//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var s_tee = "|--"
var s_bot = "`--"

func normpath(path []string) string {
	name := strings.Join(path, "/")
	if len(name) > 2 && name[:2] == "//" {
		name = name[1:]
	}
	return name
}

func inspect(dir *groot.Directory, path []string, indent string) {
	name := normpath(path)
	if dir == nil {
		fmt.Printf("err: invalid directory [%s]\n", name)
		return
	}
	keys := dir.Keys()
	nkeys := len(keys)
	str := s_tee
	//fmt.Printf("%s%s -> #%d key(s)\n", indent, name, len(keys))
	for i, k := range keys {
		if i+1 >= nkeys {
			str = s_bot
		}
		switch v := k.Value().(type) {
		default:
			fmt.Printf("%s%s %s title='%s' type=%s\n",
				indent, str, k.Name(), k.Title(), k.Class())
			
		case *groot.Directory:
			fmt.Printf("%s%s %s title='%s' type=%s\n",
				indent, str, 
				k.Name(), k.Title(), k.Class())
			path := append(path, k.Name())
			inspect(v, path, indent+"    ")

		case *groot.Tree:
			nbranches := len(v.Branches())
			fmt.Printf("%s%s %s title='%s' entries=%v nbranches=%v type=%s\n",
				indent, str, 
				k.Name(), k.Title(), v.Entries(), nbranches, k.Class())
			if *detailed {
				strbr := s_tee
				for i,branch := range v.Branches() {
					if i+1 >= nbranches {
						strbr = s_bot
					}
					fmt.Printf(" %s%s%s %s type=%s\n",
						indent, "   ", strbr, branch.Name(), branch.Class())
				}
			}
		}
	}
}

func main() {
	fmt.Printf(":: groot-ls ::\n")
	flag.Parse()

	// if *cpuprofile != "" {
    //     f, err := os.Create(*cpuprofile)
    //     if err != nil {
    //         log.Fatal(err)
    //     }
    //     pprof.StartCPUProfile(f)
    //     defer pprof.StopCPUProfile()
    // }

	if *fname == "" {
		fmt.Printf("**error** you have to give a (valid) path to a ROOT file\n")
		os.Exit(1)
	}

	f, err := groot.NewFileReader(*fname)
	if err != nil {
		fmt.Printf("**error** %v\n", err)
		os.Exit(1)
	}

	if f == nil {
		fmt.Printf("**error** invalid file pointer\n")
		os.Exit(1)
	}

	fmt.Printf("file: '%s' (version=%v)\n", f.Name(), f.Version())

	dir := f.Dir()
	inspect(dir, []string{"/"}, "")

	fmt.Printf("::bye.\n")
}

// EOF
