package version

import (
	"fmt"
	"runtime"
)

const entryFormat = "%-8s %s\n"

var appVer = "dev"
var appArch = runtime.GOARCH

func PrintVersion() {
	fmt.Println("GHFS: Go HTTP File Server")
	fmt.Printf(entryFormat, "Version:", appVer)
	fmt.Printf(entryFormat, "SDK:", runtime.Version())
	fmt.Printf(entryFormat, "OS:", runtime.GOOS)
	fmt.Printf(entryFormat, "ARCH:", appArch)
}
