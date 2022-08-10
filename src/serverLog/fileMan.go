// fileMan maintains opened files shared by multiple logMan

package serverLog

import (
	"errors"
	"os"
	"sync"
)

const chanBuffer = 31
const fileMode = 0660
const logEnding = '\n'

type fileEntry struct {
	fsPath string
	info   os.FileInfo
	file   *os.File
	ch     chan []byte
}

type FileMan struct {
	wg      *sync.WaitGroup
	entries []*fileEntry
}

func openLogFile(fsPath string) (*os.File, error) {
	return os.OpenFile(fsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
}

func newFileEntry(fsPath string, info os.FileInfo, file *os.File) *fileEntry {
	ch := make(chan []byte, chanBuffer)

	entry := &fileEntry{
		fsPath: fsPath,
		info:   info,
		file:   file,
		ch:     ch,
	}

	return entry
}

func (entry *fileEntry) serve() {
	for payload := range entry.ch {
		payload = append(payload, logEnding)
		_, e := entry.file.Write(payload)
		if e != nil {
			os.Stderr.WriteString(e.Error() + "\n")
		}
	}
	if entry.info != nil { // not Stdout or Stderr
		entry.file.Close()
	}
}

func (entry *fileEntry) reopen() error {
	if entry.info == nil { // Stdout or Stderr
		return nil
	}

	var err error
	var info os.FileInfo
	var file *os.File

	info, err = os.Stat(entry.fsPath)
	if os.IsNotExist(err) {
		// get file
		file, err = openLogFile(entry.fsPath)
		if err != nil {
			return err
		}
		// get info
		info, err = os.Stat(entry.fsPath)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}

		if os.SameFile(info, entry.info) {
			return nil
		}

		// use existing info, get file
		file, err = openLogFile(entry.fsPath)
		if err != nil {
			return err
		}
	}

	oldFile := entry.file
	entry.info = info
	entry.file = file

	return oldFile.Close()
}

func (entry *fileEntry) close() {
	close(entry.ch)
}

func (fMan *FileMan) getWritingCh(fsPath string, file *os.File) (chan<- []byte, error) {
	if len(fsPath) == 0 && file == nil {
		return nil, errors.New("log file not provided")
	}

	var err error
	var info os.FileInfo

	if file == nil { // regular file
		info, err = os.Stat(fsPath)
		if os.IsNotExist(err) {
			// get file
			file, err = openLogFile(fsPath)
			if err != nil {
				return nil, err
			}
			// get info
			info, err = os.Stat(fsPath)
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil {
				return nil, err
			}

			for _, entry := range fMan.entries {
				if os.SameFile(info, entry.info) {
					return entry.ch, nil
				}
			}

			// use existing info, get file
			file, err = openLogFile(fsPath)
			if err != nil {
				return nil, err
			}
		}
	} else { // Stdout or Stderr
		for _, entry := range fMan.entries {
			if file == entry.file {
				return entry.ch, nil
			}
		}

		fsPath = file.Name()
	}

	entry := newFileEntry(fsPath, info, file)
	fMan.entries = append(fMan.entries, entry)

	fMan.wg.Add(1)
	go func() {
		entry.serve()
		fMan.wg.Done()
	}()

	return entry.ch, nil
}

func (fMan *FileMan) Reopen() []error {
	var errs []error

	for _, entry := range fMan.entries {
		err := entry.reopen()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (fMan *FileMan) Close() {
	for _, entry := range fMan.entries {
		entry.close()
	}
	fMan.wg.Wait()
}

func (fMan *FileMan) NewLogger(accLogFilename, errLogFilename string) (*Logger, []error) {
	return newLogger(fMan, accLogFilename, errLogFilename)
}

func NewFileMan() *FileMan {
	return &FileMan{
		wg: &sync.WaitGroup{},
	}
}
