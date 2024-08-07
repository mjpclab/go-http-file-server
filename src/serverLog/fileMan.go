package serverLog

import (
	"errors"
	"os"
	"sync"
)

const fileMode = 0660

type FileMan struct {
	wg    *sync.WaitGroup
	dests []*fileDest
}

func (fMan *FileMan) ReOpen() []error {
	var errs []error

	for _, dest := range fMan.dests {
		err := dest.reopen()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (fMan *FileMan) Close() {
	for _, dest := range fMan.dests {
		dest.close()
	}
	fMan.wg.Wait()
}

func (fMan *FileMan) getWritingCh(fsPath string) (chan<- []byte, error) {
	if len(fsPath) == 0 {
		return nil, errors.New("log file not provided")
	}

	var file *os.File
	var info os.FileInfo
	var err error

	var ch chan<- []byte
	file, info, err = getFileInfoIfNotMatch(fsPath, func(info os.FileInfo) bool {
		for _, dest := range fMan.dests {
			if os.SameFile(info, dest.info) {
				ch = dest.ch
				return true
			}
		}
		return false
	})

	if err != nil {
		return nil, err
	}

	if ch != nil {
		return ch, nil
	}

	dest := newFileDest(fsPath, file, info)
	fMan.dests = append(fMan.dests, dest)

	fMan.wg.Add(1)
	go func() {
		dest.serve()
		fMan.wg.Done()
	}()

	return dest.ch, nil
}

func (fMan *FileMan) newLogChan(fsPath string) (loggerChan, error) {
	var ch chan<- []byte
	var err error

	if len(fsPath) > 0 {
		ch, err = fMan.getWritingCh(fsPath)
		if err != nil {
			return nil, err
		}
	}

	return ch, nil
}

func (fMan *FileMan) NewLogger(accLogFilename, errLogFilename string) (*Logger, []error) {
	var errs []error

	accChan, err := fMan.newLogChan(accLogFilename)
	if err != nil {
		errs = append(errs, err)
	}

	errChan, err := fMan.newLogChan(errLogFilename)
	if err != nil {
		errs = append(errs, err)
	}

	logger := &Logger{
		acc: accChan,
		err: errChan,
	}
	return logger, nil
}

func NewFileMan() *FileMan {
	return &FileMan{
		wg: &sync.WaitGroup{},
	}
}
