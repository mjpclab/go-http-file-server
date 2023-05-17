//go:build !windows
// +build !windows

package util

var IsPathEqual = IsStrEqualAccurate

var HasUrlPrefixDir = HasUrlPrefixDirAccurate
var HasFsPrefixDir = HasFsPrefixDirAccurate
