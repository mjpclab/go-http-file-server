package serverLog

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestMatchOrOpenFile(t *testing.T) {
	var err error

	f, err := os.CreateTemp("", "log-*")
	if err != nil {
		t.Fatal("create temp log file failed")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	i, err := f.Stat()
	if err != nil {
		t.Fatal("stat temp log file failed")
	}

	// matched, returns nil
	file, info, err := matchOrOpenFile(f.Name(), func(_info os.FileInfo) bool {
		if !os.SameFile(i, _info) {
			t.Error("callback `_info` should point to the same file as `i`")
			return false
		} else {
			return true
		}
	})
	if file != nil {
		t.Error("`file` should be nil since match func returns true")
		file.Close()
	}
	if info != nil {
		t.Error("`info` should be nil since match func returns true")
	}

	// not matched, returns file/info
	file, info, err = matchOrOpenFile(f.Name(), func(_info os.FileInfo) bool {
		return false
	})
	if file == nil {
		t.Error()
	} else {
		file.Close()
	}
	if info == nil {
		t.Error()
	} else if !os.SameFile(i, info) {
		t.Error(i.Name(), info.Name())
	}

	// target is directory, returns error
	file, info, err = matchOrOpenFile(os.TempDir(), func(_info os.FileInfo) bool {
		return false
	})
	if file != nil || info != nil || err == nil {
		t.Error()
	}

	// create file and returns file/info if not exists
	file, info, err = matchOrOpenFile(f.Name()+"-1", func(_info os.FileInfo) bool {
		return os.SameFile(i, _info)
	})
	if file == nil {
		t.Error()
	} else {
		file.Close()
		os.Remove(file.Name())
	}
	if info == nil {
		t.Error()
	} else if os.SameFile(i, info) {
		t.Error(i.Name(), info.Name())
	}
}

func TestFileManNewLogChan(t *testing.T) {
	var err error

	file1, err := os.CreateTemp("", "log-*")
	if err != nil {
		t.Fatal("create log file 1 failed")
	}
	file1.Close()
	os.Remove(file1.Name())

	file2, err := os.CreateTemp("", "log-*")
	if err != nil {
		t.Fatal("create log file 2 failed")
	}
	file2.Close()
	os.Remove(file2.Name())

	man := NewFileMan()
	ch1a, _ := man.newLogChan(file1.Name())
	ch1b, _ := man.newLogChan(file1.Name())
	if ch1a != ch1b {
		t.Error()
	}

	ch2, _ := man.newLogChan(file2.Name())
	if ch2 == ch1a {
		t.Error()
	}
	man.Close()
}

func TestFileMan(t *testing.T) {
	var err error

	accLogFile1, err := os.CreateTemp("", "log-*")
	if err != nil {
		t.Fatal("create log file 1 failed")
	}
	accLogFile1.Close()
	os.Remove(accLogFile1.Name())

	accLogFile2, err := os.CreateTemp("", "log-*")
	if err != nil {
		t.Fatal("create log file 2 failed")
	}
	accLogFile2.Close()
	os.Remove(accLogFile2.Name())

	errLogFile, err := os.CreateTemp("", "log-*")
	if err != nil {
		t.Fatal("create error log file failed")
	}
	errLogFile.Close()
	os.Remove(errLogFile.Name())

	var es []error
	man := NewFileMan()
	logger1, es := man.NewLogger(accLogFile1.Name(), errLogFile.Name())
	if len(es) > 0 {
		t.Fatal(es)
	}
	logger2, es := man.NewLogger(accLogFile2.Name(), errLogFile.Name())
	if len(es) > 0 {
		t.Fatal(es)
	}

	logger1.LogAccessString("host 1 access allowed")
	logger1.LogErrorString("host 1 access denied")
	logger2.LogAccess([]byte("host 2 access allowed"))
	logger2.LogError([]byte("host 2 access denied"))

	time.Sleep(100 * time.Millisecond) // wait for log written
	man.ReOpen()
	logger1.LogAccessString("HOST 1 access allowed")
	logger1.LogErrorString("HOST 1 access denied")
	logger2.LogAccessString("HOST 2 access allowed")
	logger2.LogErrorString("HOST 2 access denied")

	err = os.Rename(errLogFile.Name(), errLogFile.Name()+"-old")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond) // wait for log written
	man.ReOpen()

	logger1.LogErrorString("403 HOST I ACCESS DENIED")
	logger2.LogErrorString("403 HOST II ACCESS DENIED")
	man.Close()

	accLog1, _ := os.ReadFile(accLogFile1.Name())
	if !bytes.Equal(accLog1, []byte("host 1 access allowed\nHOST 1 access allowed\n")) {
		t.Error(string(accLog1))
	}
	accLog2, _ := os.ReadFile(accLogFile2.Name())
	if !bytes.Equal(accLog2, []byte("host 2 access allowed\nHOST 2 access allowed\n")) {
		t.Error(string(accLog2))
	}

	errLogOld, _ := os.ReadFile(errLogFile.Name() + "-old")
	if !bytes.Equal(errLogOld, []byte("host 1 access denied\nhost 2 access denied\nHOST 1 access denied\nHOST 2 access denied\n")) {
		t.Error(string(errLogOld))
	}
	errLogNew, _ := os.ReadFile(errLogFile.Name())
	if !bytes.Equal(errLogNew, []byte("403 HOST I ACCESS DENIED\n403 HOST II ACCESS DENIED\n")) {
		t.Error(string(errLogNew))
	}
}
