package serverLog

import "os"

const stdIO string = "-"

type HybridMan struct {
	fMan *FileMan
	wMan *WriterMan
}

func (man *HybridMan) ReOpen() []error {
	return man.fMan.ReOpen()
}

func (man *HybridMan) Close() {
	man.fMan.Close()
	man.wMan.Close()
}

func (man *HybridMan) NewLogger(accessLog, errorLog string) (*Logger, []error) {
	var accChan, errChan loggerChan
	var err error
	var errs []error

	if accessLog == stdIO {
		accChan = man.wMan.newLogChan(os.Stdout)
	} else {
		accChan, err = man.fMan.newLogChan(accessLog)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if errorLog == stdIO {
		errChan = man.wMan.newLogChan(os.Stderr)
	} else {
		errChan, err = man.fMan.newLogChan(errorLog)
		if err != nil {
			errs = append(errs, err)
		}
	}

	logger := &Logger{
		acc: accChan,
		err: errChan,
	}
	return logger, errs
}

func NewHybridMan() *HybridMan {
	return &HybridMan{
		fMan: NewFileMan(),
		wMan: NewWriterMan(),
	}
}
