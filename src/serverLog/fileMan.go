package serverLog

import (
	"errors"
	"os"
	"sync"
)

const fileMode = 0660

func openLogFile(fsPath string) (*os.File, error) {
	return os.OpenFile(fsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
}

func matchOrOpenFile(fsPath string, match func(info os.FileInfo) bool) (file *os.File, info os.FileInfo, err error) {
	info, err = os.Stat(fsPath)

	if os.IsNotExist(err) {
		file, err = openLogFile(fsPath)
		if err != nil {
			return nil, nil, err
		}

		info, err = file.Stat()
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
		return nil, nil, nil
	}

	// use existing info, get file
	file, err = openLogFile(fsPath)
	if err != nil {
		return nil, nil, err
	}

	return
}

type FileMan struct {
	mu    sync.Mutex
	wg    *sync.WaitGroup
	dests []*fileDest
}

func (fMan *FileMan) ReOpen() []error {
	var errs []error

	fMan.mu.Lock()
	for _, dest := range fMan.dests {
		file, info, err := matchOrOpenFile(dest.fsPath, func(info os.FileInfo) bool {
			return os.SameFile(info, dest.info)
		})
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if file != nil && info != nil {
			dest.fileCh <- file
			dest.info = info
		}
	}
	fMan.mu.Unlock()

	return errs
}

func (fMan *FileMan) Close() {
	fMan.mu.Lock()
	dests := fMan.dests
	fMan.dests = nil
	fMan.mu.Unlock()

	for _, dest := range dests {
		dest.close()
	}
	fMan.wg.Wait()
}

func (fMan *FileMan) serve(dest *fileDest, file *os.File) {
	fMan.wg.Add(1)
	go func() {
		dest.serve(file)
		fMan.wg.Done()
	}()
}

func (fMan *FileMan) getWritingCh(fsPath string) (chan<- []byte, error) {
	if len(fsPath) == 0 {
		return nil, errors.New("log file not provided")
	}

	var file *os.File
	var info os.FileInfo
	var err error

	var ch chan<- []byte
	fMan.mu.Lock()
	defer fMan.mu.Unlock()
	file, info, err = matchOrOpenFile(fsPath, func(info os.FileInfo) bool {
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

	dest := newFileDest(fsPath, info)
	fMan.dests = append(fMan.dests, dest)
	fMan.serve(dest, file)

	return dest.ch, nil
}

func (fMan *FileMan) newLogChan(fsPath string) (loggerChan, error) {
	if len(fsPath) == 0 {
		return nil, nil
	}
	return fMan.getWritingCh(fsPath)
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
	return logger, errs
}

func NewFileMan() *FileMan {
	return &FileMan{
		wg: &sync.WaitGroup{},
	}
}
