//go:build windows
// +build windows

package util

var IsPathEqual = IsStrEqualNoCase

var HasUrlPrefixDir = HasUrlPrefixDirNoCase
var HasFsPrefixDir = HasFsPrefixDirNoCase
