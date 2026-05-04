package serverLog

type loggerChan chan<- []byte

func (ch loggerChan) canLog() bool {
	return ch != nil
}

func (ch loggerChan) log(payload []byte) {
	// `payload` MUST NOT be nil, otherwise will stop the receiver of the `ch`
	if len(payload) > 0 && ch.canLog() {
		ch <- payload
	}
}
func (ch loggerChan) logString(payload string) {
	ch.log([]byte(payload))
}
