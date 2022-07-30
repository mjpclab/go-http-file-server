package serverError

import (
	"fmt"
	"os"
)

func AppendError(errs []error, err error) []error {
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

func CheckError(errs ...error) bool {
	hasError := false

	for _, err := range errs {
		if err == nil {
			continue
		}
		hasError = true
		fmt.Fprintln(os.Stderr, err)
	}

	return hasError
}

func CheckFatal(errs ...error) bool {
	hasError := CheckError(errs...)

	if hasError {
		os.Exit(1)
	}

	return hasError
}
