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
