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

func getLogEntry(payload string) string {
	buffer := &bytes.Buffer{}
	buffer.WriteString(util.FormatTimeNanosecond(time.Now()))
	buffer.WriteByte(' ')
	buffer.WriteString(payload)
	buffer.WriteByte('\n')

	return buffer.String()
}

func (l *Logger) CanLogAccess() bool {
	return l.accLogFile != nil
}

func (l *Logger) CanLogError() bool {
	return l.errLogFile != nil
}

func (l *Logger) LogAccess(payload string) {
	if !l.CanLogAccess() {
		return
	}

	_, e := l.accLogFile.WriteString(getLogEntry(payload))
	if e != nil {
		l.LogError(e.Error())
	}
}

func (l *Logger) LogError(payload string) {
	if !l.CanLogError() {
		return
	}

	_, e := l.errLogFile.WriteString(getLogEntry(payload))
	if e != nil {
		os.Stdout.WriteString(e.Error() + "\n")
	}
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
