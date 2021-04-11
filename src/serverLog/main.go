package serverLog

import (
	"os"
)

const CHAN_BUFFER = 15

type Logger struct {
	accLogMan *logMan
	errLogMan *logMan
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

func (l *Logger) LogErrors(errors ...error) {
	for _, err := range errors {
		l.errLogMan.LogString(err.Error())
	}
}

func (l *Logger) LogErrorString(payload string) {
	l.errLogMan.LogString(payload)
}

func (l *Logger) Open() (errors []error) {
	if err := l.accLogMan.Open(); err != nil {
		errors = append(errors, err)
	} else {
		l.accLogMan.Enable()
	}

	if err := l.errLogMan.Open(); err != nil {
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
	if err := l.accLogMan.ReOpen(); err != nil {
		errors = append(errors, err)
	}

	if err := l.errLogMan.ReOpen(); err != nil {
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
