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

}

// EOF
