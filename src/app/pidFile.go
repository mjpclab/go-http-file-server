package app

import (
	"mjpclab.dev/ghfs/src/serverError"
	"os"
	"strconv"
)

func writePidFile(pidFilePath string) (errs []error) {
	pidContent := strconv.Itoa(os.Getpid())
	pidFile, err := os.OpenFile(pidFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return []error{err}
	}

	_, err = pidFile.WriteString(pidContent)
	errs = serverError.AppendError(errs, err)

	err = pidFile.Close()
	errs = serverError.AppendError(errs, err)

	return
}
