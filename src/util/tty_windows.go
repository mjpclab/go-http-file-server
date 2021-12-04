//go:build windows
// +build windows

package util

import "os"

func GetTTYFile() (*os.File, func() error) {
	return os.Stdout, noopNilError
}
