// logMan controls if log action should really be performed by `canLog()`
// It also provides utility functions to accept other data format rather than []byte,
// e.g. `logString()` accepts string as log data.

package serverLog

import "os"

type logMan chan<- []byte

func (lMan logMan) canLog() bool {
	return lMan != nil
}

func (lMan logMan) log(payload []byte) {
	if lMan.canLog() {
		lMan <- payload
	}
}
func (lMan logMan) logString(payload string) {
	lMan.log([]byte(payload))
}

func newLogMan(fMan *FileMan, fsPath string, dashFile *os.File) (logMan, error) {
	var ch chan<- []byte
	var err error

	if len(fsPath) > 0 {
		if fsPath == "-" {
			ch, err = fMan.getWritingCh("", dashFile)
		} else {
			ch, err = fMan.getWritingCh(fsPath, nil)
		}

		if err != nil {
			return nil, err
		}
	}

	man := logMan(ch)
	return man, nil
}
