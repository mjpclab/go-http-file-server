package shimgo

import (
	"syscall"
)

func Bytes_LastIndexByte(s []byte, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func Os_LookupEnv(key string) (string, bool) {
	return syscall.Getenv(key)
}
