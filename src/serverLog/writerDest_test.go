package serverLog

import (
	"bytes"
	"testing"
)

func TestWriterDest(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	dest := newWriterDest(buf)

	dest.ch <- []byte("hello")
	dest.ch <- []byte("world")

	if buf.Len() != 0 {
		t.Error(buf.Len())
	}

	go dest.close()
	dest.serve()
	if buf.String() != "hello\nworld\n" {
		t.Error(buf.String())
	}
}
