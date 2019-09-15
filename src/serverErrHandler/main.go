package serverErrHandler

import (
	"../serverLog"
	"errors"
	"fmt"
	"os"
)

type ErrHandler struct {
	logger *serverLog.Logger
}

func NewErrHandler(logger *serverLog.Logger) *ErrHandler {
	return &ErrHandler{logger}
}

func (h *ErrHandler) LogError(err error) bool {
	if err == nil {
		return false
	}
	if h.logger == nil {
		CheckError(errors.New("logger not initialized"))
		CheckError(err)
		return true
	}

	h.logger.LogErrorString(err.Error())
	return true
}

func (h *ErrHandler) LogFatal(err error) bool {
	hasError := h.LogError(err)

	if hasError {
		os.Exit(1)
	}

	return hasError
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
