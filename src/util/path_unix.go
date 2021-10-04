//go:build !windows
// +build !windows

package util

import "path/filepath"

func NormalizeFsPath(input string) (string, error) {
	return filepath.Abs(input)
}
