package serverLog

type Logger struct {
	acc loggerChan
	err loggerChan
}

func (l *Logger) CanLogAccess() bool {
	return l.acc.canLog()
}

func (l *Logger) CanLogError() bool {
	return l.err.canLog()
}

func (l *Logger) LogAccess(payload []byte) {
	l.acc.log(payload)
}
func (l *Logger) LogAccessString(payload string) {
	l.acc.logString(payload)
}

func (l *Logger) LogError(payload []byte) {
	l.err.log(payload)
}

func (l *Logger) LogErrorString(payload string) {
	l.err.logString(payload)
}

func (l *Logger) LogErrors(errors ...error) {
	for _, err := range errors {
		l.err.logString(err.Error())
	}
}
