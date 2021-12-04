//go:build !windows
// +build !windows

package util

import "os"

func GetTTYFile() (file *os.File, teardown func() error) {
	tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0220)
	if err == nil {
		return tty, tty.Close
	}

	return os.Stdout, noopNilError
}
