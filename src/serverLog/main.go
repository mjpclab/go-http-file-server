package serverLog

import (
	"../util"
	"bytes"
	"os"
	"sync"
	"time"
)

const CHAN_BUFFER = 7

type Logger struct {
	accLogFile *os.File
	accLogChan chan []byte

	errLogFile *os.File
	errLogChan chan []byte

	waitGroup sync.WaitGroup
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
	return l.accLogFile != nil
}

func (l *Logger) CanLogError() bool {
	return l.errLogFile != nil
}

func (l *Logger) LogAccess(payload []byte) {
	if l.CanLogAccess() {
		l.accLogChan <- payload
	}
}
func (l *Logger) LogAccessString(payload string) {
	l.LogAccess([]byte(payload))
}

func (l *Logger) LogError(payload []byte) {
	if l.CanLogError() {
		l.errLogChan <- payload
	}
}

func (l *Logger) LogErrorString(payload string) {
	l.LogError([]byte(payload))
}

func (l *Logger) enableAccLog() {
	for {
		payload, ok := <-l.accLogChan
		if !ok {
			break
		}

		_, e := l.accLogFile.Write(getLogEntry(payload))
		if e != nil {
			l.LogError([]byte(e.Error()))
		}
	}
	l.waitGroup.Done()
}

func (l *Logger) enableErrLog() {
	for {
		payload, ok := <-l.errLogChan
		if !ok {
			break
		}

		_, e := l.errLogFile.Write(getLogEntry(payload))
		if e != nil {
			os.Stdout.WriteString(e.Error() + "\n")
		}
	}
	l.waitGroup.Done()
}

func (l *Logger) Close() {
	if l.accLogChan != nil {
		close(l.accLogChan)
		l.accLogChan = nil
	}
	if l.errLogChan != nil {
		close(l.errLogChan)
		l.errLogChan = nil
	}

	l.waitGroup.Wait()

	if l.accLogFile != nil {
		l.accLogFile.Close()
		l.accLogFile = nil
	}

	if l.errLogFile != nil {
		l.errLogFile.Close()
		l.errLogFile = nil
	}
}

func NewLogger(accessFilename, errorFilename string) (*Logger, error) {
	var accLogFile, errLogFile *os.File
	var accLogChan, errLogChan chan []byte

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
		accLogChan = make(chan []byte, CHAN_BUFFER)
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
		errLogChan = make(chan []byte, CHAN_BUFFER)
	}

	logger := &Logger{
		accLogFile: accLogFile,
		accLogChan: accLogChan,

		errLogFile: errLogFile,
		errLogChan: errLogChan,
	}

	if logger.accLogChan != nil {
		logger.waitGroup.Add(1)
		go logger.enableAccLog()
	}
	if logger.errLogChan != nil {
		logger.waitGroup.Add(1)
		go logger.enableErrLog()
	}

	return logger, nil
}
