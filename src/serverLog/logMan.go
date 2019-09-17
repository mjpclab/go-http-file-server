package serverLog

import (
	"os"
	"sync"
)

type logMan struct {
	filename string
	file     *os.File
	ch       chan []byte

	dashFile *os.File
	wg       sync.WaitGroup
}

func (l *logMan) Open() (err error) {
	if len(l.filename) == 0 {
		return
	}

	if l.filename == "-" {
		l.file = l.dashFile
	} else {
		l.file, err = os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return
		}
	}

	return
}

func (l *logMan) Enable() {
	if l.file == nil {
		return
	}

	ch := make(chan []byte, CHAN_BUFFER)
	l.ch = ch

	l.wg.Add(1)
	go func() {
		for {
			payload, ok := <-ch
			if !ok {
				break
			}

			_, e := l.file.Write(getLogEntry(payload))
			if e != nil {
				os.Stderr.WriteString(e.Error() + "\n")
			}
		}
		l.wg.Done()
	}()
}

func (l *logMan) Close() {
	if l.ch != nil {
		close(l.ch)
		l.ch = nil
	}

	l.wg.Wait()

	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

func (l *logMan) ReOpen() (err error) {
	oldFile := l.file
	if oldFile == nil || oldFile == l.dashFile {
		return
	}

	err = l.Open()
	if err != nil {
		return
	}

	err = oldFile.Close()
	return
}

func (l *logMan) CanLog() bool {
	return l.file != nil && l.ch != nil
}

func (l *logMan) Log(payload []byte) {
	if l.CanLog() {
		l.ch <- payload
	}
}
func (l *logMan) LogString(payload string) {
	l.Log([]byte(payload))
}

func NewLogMan(filename string, dashFile *os.File) *logMan {
	return &logMan{
		filename: filename,
		dashFile: dashFile,
	}
}
