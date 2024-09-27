package serverLog

import (
	"io"
	"sync"
)

type WriterMan struct {
	wg    *sync.WaitGroup
	dests []*writerDest
}

func (wMan *WriterMan) Close() {
	for _, dest := range wMan.dests {
		dest.close()
	}
	wMan.wg.Wait()
}

func (wMan *WriterMan) getWritingCh(w io.Writer) (chan<- []byte, error) {
	for _, dest := range wMan.dests {
		if w == dest.w {
			return dest.ch, nil
		}
	}

	dest := newWriterDest(w)
	wMan.dests = append(wMan.dests, dest)

	wMan.wg.Add(1)
	go func() {
		dest.serve()
		wMan.wg.Done()
	}()

	return dest.ch, nil
}

func (wMan *WriterMan) newLogChan(w io.Writer) (loggerChan, error) {
	var ch chan<- []byte
	var err error

	if w != nil {
		ch, err = wMan.getWritingCh(w)
		if err != nil {
			return nil, err
		}
	}

	return ch, nil
}

func (wMan *WriterMan) NewLogger(accLogWriter, errLogWriter io.Writer) (*Logger, []error) {
	var errs []error

	accChan, err := wMan.newLogChan(accLogWriter)
	if err != nil {
		errs = append(errs, err)
	}

	errChan, err := wMan.newLogChan(errLogWriter)
	if err != nil {
		errs = append(errs, err)
	}

	logger := &Logger{
		acc: accChan,
		err: errChan,
	}
	return logger, nil
}

func NewWriterMan() *WriterMan {
	return &WriterMan{
		wg: &sync.WaitGroup{},
	}
}
