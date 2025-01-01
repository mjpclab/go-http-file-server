package serverLog

import (
	"os"
)

type fileDest struct {
	fsPath string
	file   *os.File
	info   os.FileInfo
	ch     chan []byte
}

func newFileDest(fsPath string, file *os.File, info os.FileInfo) *fileDest {
	ch := make(chan []byte, logQueueSize)

	dest := &fileDest{
		fsPath: fsPath,
		file:   file,
		info:   info,
		ch:     ch,
	}

	return dest
}

func (dest *fileDest) serve() {
	for payload := range dest.ch {
		payload = append(payload, logEnding)
		_, e := dest.file.Write(payload)
		if e != nil {
			os.Stderr.WriteString(e.Error() + "\n")
		}
	}
	dest.file.Close()
}

func (dest *fileDest) close() {
	close(dest.ch)
}
