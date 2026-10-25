package serverLog

type Logger struct {
	acc loggerChan
	mut loggerChan
	err loggerChan
}

func (l *Logger) CanLogAccess() bool {
	return l.acc.canLog()
}

func (l *Logger) CanLogMutation() bool {
	return l.mut.canLog()
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

func (l *Logger) LogMutation(payload []byte) {
	l.mut.log(payload)
}

func (l *Logger) LogMutationString(payload string) {
	l.mut.logString(payload)
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
