//go:build !windows
// +build !windows

package util

import "path/filepath"

func IsPathEqual(a, b string) bool {
	return a == b
}

func HasUrlPrefixDir(urlPath, prefix string) bool {
	return hasPrefixDirAccurate(urlPath, prefix, '/')
}

func HasFsPrefixDir(fsPath, prefix string) bool {
	return hasPrefixDirAccurate(fsPath, prefix, filepath.Separator)
}
