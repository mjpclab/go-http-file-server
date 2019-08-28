package serverError

import (
	"../serverLog"
	"errors"
	"fmt"
	"os"
)

var logger *serverLog.Logger

func SetLogger(instance *serverLog.Logger) {
	logger = instance
}

func LogError(err error) bool {
	if err == nil {
		return false
	}
	if logger == nil {
		CheckError(errors.New("logger not initialized"))
		CheckError(err)
		return true
	}

	logger.LogError(err.Error())
	return true
}

func LogFatal(err error) bool {
	hasError := LogError(err)

	if hasError {
		os.Exit(1)
	}

	return hasError
}

func CheckError(err error) bool {
	if err == nil {
		return false
	}

	fmt.Fprintln(os.Stderr, err)
	return true
}

func CheckFatal(err error) bool {
	hasError := CheckError(err)

	if hasError {
		os.Exit(1)
	}

	return hasError
}
