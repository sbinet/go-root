package main

import (
	"fmt"
	"os"

	"bitbucket.org/binet/go-root/pkg/groot"
)

func main() {
	fmt.Printf("== test go-root ==\n")
	f, err := groot.NewFileReader("ntuple.0.root")
	if err != nil {
		fmt.Printf("**error**: %v\n", err)
		os.Exit(1)
	}

	if f == nil {
		fmt.Printf("invalid file pointer\n")
		os.Exit(1)
	}

	fmt.Printf("f: %s (version=%v)\n", f.GetName(), f.GetVersion())

}

// EOF
