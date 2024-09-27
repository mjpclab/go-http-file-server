package serverLog

import (
	"bytes"
	"testing"
)

func TestWriterManNewLogChan(t *testing.T) {
	buf1 := bytes.NewBuffer(make([]byte, 0, 128))
	buf2 := bytes.NewBuffer(make([]byte, 0, 128))

	man := NewWriterMan()
	ch1a, _ := man.newLogChan(buf1)
	ch1b, _ := man.newLogChan(buf1)
	if ch1a != ch1b {
		t.Error()
	}

	ch2, _ := man.newLogChan(buf2)
	if ch2 == ch1a {
		t.Error()
	}
}

func TestWriterManNewLogger(t *testing.T) {
	acc1File := bytes.NewBuffer(make([]byte, 0, 128))
	acc2File := bytes.NewBuffer(make([]byte, 0, 128))
	errFile := bytes.NewBuffer(make([]byte, 0, 128))

	man := NewWriterMan()

	logger1, _ := man.NewLogger(acc1File, errFile)
	logger1.LogAccessString("host 1 access allowed")
	logger1.LogErrorString("host 1 access denied")

	logger2, _ := man.NewLogger(acc2File, errFile)
	logger2.LogAccessString("host 2 access allowed")
	logger2.LogErrorString("host 2 access denied")

	man.Close()

	if acc1File.String() != "host 1 access allowed\n" {
		t.Error()
	}
	if acc2File.String() != "host 2 access allowed\n" {
		t.Error()
	}
	if errFile.String() != "host 1 access denied\nhost 2 access denied\n" {
		t.Error()
	}
}
