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
}

func (man *HybridMan) NewLogger(accessLog, errorLog string) (*Logger, []error) {
	var accChan, errChan loggerChan
	var err error
	var errs []error

	if accessLog != stdIO {
		accChan, err = man.fMan.newLogChan(accessLog)
	} else {
		accChan, err = man.wMan.newLogChan(os.Stdout)
	}
	if err != nil {
		errs = append(errs, err)
	}

	if errorLog != stdIO {
		errChan, err = man.fMan.newLogChan(errorLog)
	} else {
		errChan, err = man.wMan.newLogChan(os.Stderr)
	}
	if err != nil {
		errs = append(errs, err)
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
