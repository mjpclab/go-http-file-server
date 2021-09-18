//go:build windows
// +build windows

package util

import "path/filepath"

func NormalizeFsPath(input string) (string, error) {
	abs, err := filepath.Abs(input)
	if err != nil {
		return abs, err
	}

	abs = AsciiToLowerCase(abs)

	return abs, err
}
