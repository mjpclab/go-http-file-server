package version

import (
	"fmt"
	"runtime"
)

var appVer = "dev"

const entryFormat = "%-8s %s\n"

func PrintVersion() {
	fmt.Println("GHFS: Go HTTP File Server")
	fmt.Printf(entryFormat, "Version:", appVer)
	fmt.Printf(entryFormat, "SDK:", runtime.Version())
	fmt.Printf(entryFormat, "OS:", runtime.GOOS)
	fmt.Printf(entryFormat, "ARCH:", runtime.GOARCH)
}
