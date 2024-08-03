package serverLog

import (
	"errors"
	"os"
)

func openLogFile(fsPath string) (*os.File, error) {
	return os.OpenFile(fsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
}

func getFileInfoIfNotMatch(fsPath string, match func(info os.FileInfo) bool) (file *os.File, info os.FileInfo, err error) {
	info, err = os.Stat(fsPath)

	if os.IsNotExist(err) {
		file, err = openLogFile(fsPath)
		if err != nil {
			return nil, nil, err
		}

		info, err = os.Stat(fsPath)
		if err != nil {
			file.Close()
			return nil, nil, err
		}

		return
	}

	if err != nil {
		return nil, nil, err
	}

	if info.IsDir() {
		err = errors.New("should not be a directory")
		return nil, nil, err
	}

	if match(info) {
		return nil, nil, err
	}

	// use existing info, get file
	file, err = openLogFile(fsPath)
	if err != nil {
		return nil, nil, err
	}

	return
}

type fileDest struct {
	fsPath string
	file   *os.File
	info   os.FileInfo
	ch     chan []byte
}

func newFileDest(fsPath string, file *os.File, info os.FileInfo) *fileDest {
	ch := make(chan []byte, chanBuffer)

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

func (dest *fileDest) reopen() error {
	matched := false
	file, info, err := getFileInfoIfNotMatch(dest.fsPath, func(info os.FileInfo) bool {
		matched = os.SameFile(info, dest.info)
		return matched
	})

	if err != nil {
		return err
	}

	if matched {
		return nil
	}

	oldFile := dest.file
	dest.info = info
	dest.file = file
	return oldFile.Close()
}

func (dest *fileDest) close() {
	close(dest.ch)
}
