//go:build !windows
// +build !windows

package util

import "path/filepath"

var IsPathEqual = IsStrEqualAccurate

var HasUrlPrefixDir = HasUrlPrefixDirAccurate
var HasFsPrefixDir = HasFsPrefixDirAccurate

func NormalizeFsPath(input string) (string, error) {
	return filepath.Abs(input)
}
