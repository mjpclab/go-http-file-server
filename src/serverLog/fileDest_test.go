package serverLog

import (
	"bytes"
	"os"
	"testing"
)

func TestFileDest(t *testing.T) {
	file, err := os.CreateTemp("", "log-*")
	if err != nil {
		t.Fatal("create file failed")
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	info, err := file.Stat()
	if err != nil {
		t.Fatal("stat file failed")
	}

	dest := newFileDest(file.Name(), file, info)
	dest.ch <- []byte("hello")
	dest.ch <- []byte("world")

	go dest.close()
	dest.serve()

	logs, _ := os.ReadFile(file.Name())
	if !bytes.Equal(logs, []byte("hello\nworld\n")) {
		t.Error(string(logs))
	}
}
