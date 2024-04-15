package main

import (
	"mjpclab.dev/ghfs/src"
	"os"
)

func main() {
	ok := src.Main()
	if !ok {
		os.Exit(1)
	}
}
