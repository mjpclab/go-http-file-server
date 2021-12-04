// Logger maintains access log and error log objects for each virtual host.
// It also determines what a "-" file means, Stdout or Stderr, depends on the logMan's role.

package serverLog

import (
	"os"
)

type Logger struct {
	accLogMan logMan
	errLogMan logMan
}

func (l *Logger) CanLogAccess() bool {
	return l.accLogMan.canLog()
}

func (l *Logger) CanLogError() bool {
	return l.errLogMan.canLog()
}

func (l *Logger) LogAccess(payload []byte) {
	l.accLogMan.log(payload)
}
func (l *Logger) LogAccessString(payload string) {
	l.accLogMan.logString(payload)
}

func (l *Logger) LogError(payload []byte) {
	l.errLogMan.log(payload)
}

func (l *Logger) LogErrorString(payload string) {
	l.errLogMan.logString(payload)
}

func (l *Logger) LogErrors(errors ...error) {
	for _, err := range errors {
		l.errLogMan.logString(err.Error())
	}
}

func newLogger(fileMan *FileMan, accLogFilename, errLogFilename string) (*Logger, []error) {
	var errs []error

	accLogMan, err := newLogMan(fileMan, accLogFilename, os.Stdout)
	if err != nil {
		errs = append(errs, err)
	}

	errLogMan, err := newLogMan(fileMan, errLogFilename, os.Stderr)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	logger := &Logger{
		accLogMan: accLogMan,
		errLogMan: errLogMan,
	}
	return logger, nil
}
