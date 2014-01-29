package main

import (
	"github.com/fraenkel/candiedyaml"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s file1.yaml ...\n", os.Args[0])
		os.Exit(1)
	}

	candiedyaml.Run_parser(os.Args[0], os.Args[1:])
}
