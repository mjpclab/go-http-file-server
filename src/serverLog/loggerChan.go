// loggerChan controls if log action should really be performed by `canLog()`
// It also provides utility functions to accept other data format rather than []byte,
// e.g. `logString()` accepts string as log data.

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
