package serverLog

import "os"

const stdIO string = "-"

type Man struct {
	fMan *FileMan
	wMan *WriterMan
}

func (man *Man) ReOpenFiles() []error {
	return man.fMan.ReOpen()
}

func (man *Man) CloseFiles() {
	man.fMan.Close()
}

func (man *Man) NewLogger(accessLog, errorLog string) (*Logger, []error) {
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
	return logger, nil
}

func NewMan() *Man {
	return &Man{
		fMan: NewFileMan(),
		wMan: NewWriterMan(),
	}
}
