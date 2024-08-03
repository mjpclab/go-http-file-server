package serverLog

type loggerChan chan<- []byte

func (ch loggerChan) canLog() bool {
	return ch != nil
}

func (ch loggerChan) log(payload []byte) {
	if ch.canLog() {
		ch <- payload
	}
}
func (ch loggerChan) logString(payload string) {
	ch.log([]byte(payload))
}
