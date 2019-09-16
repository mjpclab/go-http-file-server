package serverLog

import (
	"../util"
	"bytes"
	"os"
	"time"
)

const CHAN_BUFFER = 7

type Logger struct {
	accLogMan *logMan
	errLogMan *logMan
}

func getLogEntry(payload []byte) []byte {
	buffer := &bytes.Buffer{}
	buffer.WriteString(util.FormatTimeSecond(time.Now()))
	buffer.WriteByte(' ')
	buffer.Write(payload)
	buffer.WriteByte('\n')

	return buffer.Bytes()
}

func (l *Logger) CanLogAccess() bool {
	return l.accLogMan.CanLog()
}

func (l *Logger) CanLogError() bool {
	return l.errLogMan.CanLog()
}

func (l *Logger) LogAccess(payload []byte) {
	l.accLogMan.Log(payload)
}
func (l *Logger) LogAccessString(payload string) {
	l.accLogMan.LogString(payload)
}

func (l *Logger) LogError(payload []byte) {
	l.errLogMan.Log(payload)
}

func (l *Logger) LogErrorString(payload string) {
	l.errLogMan.LogString(payload)
}

func (l *Logger) Open() (errors []error) {
	var err error

	err = l.accLogMan.Open()
	if err != nil {
		errors = append(errors, err)
	} else {
		l.accLogMan.Enable()
	}

	err = l.errLogMan.Open()
	if err != nil {
		errors = append(errors, err)
	} else {
		l.errLogMan.Enable()
	}

	return
}

func (l *Logger) Close() {
	l.accLogMan.Close()
	l.errLogMan.Close()
}

func (l *Logger) ReOpen() (errors []error) {
	var err error

	err = l.accLogMan.ReOpen()
	if err != nil {
		errors = append(errors, err)
	}

	err = l.errLogMan.ReOpen()
	if err != nil {
		errors = append(errors, err)
	}

	return
}

func NewLogger(accLogFilename, errLogFilename string) *Logger {
	logger := &Logger{
		accLogMan: NewLogMan(accLogFilename, os.Stdout),
		errLogMan: NewLogMan(errLogFilename, os.Stderr),
	}
	return logger
}
