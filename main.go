package main

import (
	"os"

	"mjpclab.dev/ghfs/src"
)

func main() {
	ok := src.Main()
	if !ok {
		os.Exit(1)
	}
}
