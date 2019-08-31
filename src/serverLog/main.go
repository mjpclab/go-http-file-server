package serverLog

import (
	"../util"
	"bytes"
	"os"
	"time"
)

type Logger struct {
	accLogFile *os.File
	errLogFile *os.File
}

func getLogEntry(payload []byte) []byte {
	buffer := &bytes.Buffer{}
	buffer.WriteString(util.FormatTimeNanosecond(time.Now()))
	buffer.WriteByte(' ')
	buffer.Write(payload)
	buffer.WriteByte('\n')

	return buffer.Bytes()
}

func (l *Logger) CanLogAccess() bool {
	return l.accLogFile != nil
}

func (l *Logger) CanLogError() bool {
	return l.errLogFile != nil
}

func (l *Logger) LogAccess(payload []byte) {
	if !l.CanLogAccess() {
		return
	}

	_, e := l.accLogFile.Write(getLogEntry(payload))
	if e != nil {
		l.LogError([]byte(e.Error()))
	}
}
func (l *Logger) LogAccessString(payload string) {
	l.LogAccess([]byte(payload))
}

func (l *Logger) LogError(payload []byte) {
	if !l.CanLogError() {
		return
	}

	_, e := l.errLogFile.Write(getLogEntry(payload))
	if e != nil {
		os.Stdout.WriteString(e.Error() + "\n")
	}
}

func (l *Logger) LogErrorString(payload string) {
	l.LogError([]byte(payload))
}

func NewLogger(accessFilename, errorFilename string) (*Logger, error) {
	var accLogFile, errLogFile *os.File
	var e error

	if len(accessFilename) > 0 {
		if accessFilename == "-" {
			accLogFile = os.Stdout
		} else {
			accLogFile, e = os.OpenFile(accessFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
			if e != nil {
				return nil, e
			}
		}
	}

	if len(errorFilename) > 0 {
		if errorFilename == "-" {
			errLogFile = os.Stderr
		} else {
			errLogFile, e = os.OpenFile(errorFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
			if e != nil {
				return nil, e
			}
		}
	}

	return &Logger{accLogFile, errLogFile}, nil
}
