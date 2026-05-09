package serverLog

import (
	"os"
)

type fileDest struct {
	fsPath string
	info   os.FileInfo
	ch     chan []byte
}

func newFileDest(fsPath string, info os.FileInfo) *fileDest {
	ch := make(chan []byte, logQueueSize)

	dest := &fileDest{
		fsPath: fsPath,
		info:   info,
		ch:     ch,
	}

	return dest
}

func (dest *fileDest) serve(file *os.File) {
	for payload := range dest.ch {
		if payload == nil {
			break
		}

		payload = append(payload, logEnding)
		_, e := file.Write(payload)
		if e != nil {
			os.Stderr.WriteString(e.Error() + "\n")
		}
	}
}

func (dest *fileDest) close() {
	close(dest.ch)
}
