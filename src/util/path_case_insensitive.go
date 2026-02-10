//go:build windows || darwin
// +build windows darwin

package util

import (
	"path/filepath"
	"strings"
)

func IsPathEqual(a, b string) bool {
	return strings.EqualFold(a, b)
}

func HasUrlPrefixDir(urlPath, prefix string) bool {
	return hasPrefixDirNoCase(urlPath, prefix, '/')
}

func HasFsPrefixDir(fsPath, prefix string) bool {
	return hasPrefixDirNoCase(fsPath, prefix, filepath.Separator)
}
