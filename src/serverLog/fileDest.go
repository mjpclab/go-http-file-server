package serverLog

import (
	"os"
)

type fileDest struct {
	fsPath string
	info   os.FileInfo
	ch     chan []byte
	fileCh chan *os.File
}

func newFileDest(fsPath string, info os.FileInfo) *fileDest {
	dest := &fileDest{
		fsPath: fsPath,
		info:   info,
		ch:     make(chan []byte, logQueueSize),
		fileCh: make(chan *os.File),
	}

	return dest
}

func (dest *fileDest) serve(file *os.File) {
	for {
		select {
		case payload, ok := <-dest.ch:
			if !ok {
				file.Close()
				return
			}
			payload = append(payload, logEnding)
			_, e := file.Write(payload)
			if e != nil {
				os.Stderr.WriteString(e.Error() + "\n")
			}
		case newFile := <-dest.fileCh:
			file.Close()
			file = newFile
		}
	}
}

func (dest *fileDest) close() {
	close(dest.ch)
}
