package serverLog

import (
	"io"
	"os"
)

type writerDest struct {
	w  io.Writer
	ch chan []byte
}

func newWriterDest(w io.Writer) *writerDest {
	ch := make(chan []byte, chanBuffer)

	dest := &writerDest{
		w:  w,
		ch: ch,
	}

	return dest
}

func (dest *writerDest) serve() {
	for payload := range dest.ch {
		payload = append(payload, logEnding)
		_, e := dest.w.Write(payload)
		if e != nil {
			os.Stderr.WriteString(e.Error() + "\n")
		}
	}
}

func (dest *writerDest) close() {
	close(dest.ch)
}
