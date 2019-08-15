package serverError

import (
	"fmt"
	"os"
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return true
	}

	return false
}

func CheckFatal(err error) bool {
	hasError := CheckError(err)

	if hasError {
		os.Exit(1)
	}

	return hasError
}
