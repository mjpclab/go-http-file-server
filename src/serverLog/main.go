package serverLog

import (
	"../util"
	"os"
	"strings"
	"time"
)

type Logger struct {
	accLogFile *os.File
	errLogFile *os.File
}

func getLogEntry(payload string) string {
	sb := strings.Builder{}
	sb.WriteString(util.FormatTimeNanosecond(time.Now()))
	sb.WriteByte(' ')
	sb.WriteString(payload)
	sb.WriteByte('\n')

	return sb.String()
}

func (l *Logger) AccessFileAvail() bool {
	return l.accLogFile != nil
}

func (l *Logger) ErrorFileAvail() bool {
	return l.errLogFile != nil
}

func (l *Logger) LogAccess(payload string) {
	if !l.AccessFileAvail() {
		return
	}

	_, e := l.accLogFile.WriteString(getLogEntry(payload))
	if e != nil {
		l.LogError(e.Error())
	}
}

func (l *Logger) LogError(payload string) {
	if !l.ErrorFileAvail() {
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
