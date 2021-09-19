//go:build windows
// +build windows

package util

import (
	"path/filepath"
	"strings"
)

func NormalizeFsPath(input string) (string, error) {
	abs, err := filepath.Abs(input)
	if err != nil {
		return abs, err
	}

	abs = strings.ToLower(abs)

	return abs, err
}
