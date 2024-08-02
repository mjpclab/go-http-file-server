// fileMan maintains opened files shared by multiple loggerChan

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

func (fMan *FileMan) getWritingCh(fsPath string, file *os.File) (chan<- []byte, error) {
	if len(fsPath) == 0 && file == nil {
		return nil, errors.New("log file not provided")
	}

	var info os.FileInfo
	var err error

	if file == nil { // regular file
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
	} else { // Stdout or Stderr
		for _, dest := range fMan.dests {
			if file == dest.file {
				return dest.ch, nil
			}
		}

		fsPath = file.Name()
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

func (fMan *FileMan) Reopen() []error {
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

func (fMan *FileMan) newLogChan(fsPath string, dashFile *os.File) (loggerChan, error) {
	var ch chan<- []byte
	var err error

	if len(fsPath) > 0 {
		if fsPath == "-" {
			ch, err = fMan.getWritingCh("", dashFile)
		} else {
			ch, err = fMan.getWritingCh(fsPath, nil)
		}

		if err != nil {
			return nil, err
		}
	}

	return ch, nil
}

func (fMan *FileMan) NewLogger(accLogFilename, errLogFilename string) (*Logger, []error) {
	var errs []error

	accChan, err := fMan.newLogChan(accLogFilename, os.Stdout)
	if err != nil {
		errs = append(errs, err)
	}

	errChan, err := fMan.newLogChan(errLogFilename, os.Stderr)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
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
