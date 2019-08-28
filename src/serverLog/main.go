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

func (l *Logger) Log(payload string) {
	if l.accLogFile == nil {
		return
	}

	_, e := l.accLogFile.WriteString(getLogEntry(payload))
	if e != nil {
		l.Error(e.Error())
	}
}

func (l *Logger) Error(payload string) {
	if l.errLogFile == nil {
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
