package serverLog

import (
	"io"
	"sync"
)

type WriterMan struct {
	mu    sync.Mutex
	wg    *sync.WaitGroup
	dests []*writerDest
}

func (wMan *WriterMan) ReOpen() []error {
	return nil
}

func (wMan *WriterMan) Close() {
	wMan.mu.Lock()
	dests := wMan.dests
	wMan.dests = nil
	wMan.mu.Unlock()

	for _, dest := range dests {
		dest.close()
	}
	wMan.wg.Wait()
}

func (wMan *WriterMan) getWritingCh(w io.Writer) chan<- []byte {
	wMan.mu.Lock()
	defer wMan.mu.Unlock()
	for _, dest := range wMan.dests {
		if w == dest.w {
			return dest.ch
		}
	}

	dest := newWriterDest(w)
	wMan.dests = append(wMan.dests, dest)

	wMan.wg.Add(1)
	go func() {
		dest.serve()
		wMan.wg.Done()
	}()

	return dest.ch
}

func (wMan *WriterMan) newLogChan(w io.Writer) loggerChan {
	if w == nil {
		return nil
	}
	return wMan.getWritingCh(w)
}

// NewLogger only accepts comparable Writers
func (wMan *WriterMan) NewLogger(accLogWriter, errLogWriter io.Writer) *Logger {
	accChan := wMan.newLogChan(accLogWriter)
	errChan := wMan.newLogChan(errLogWriter)
	logger := &Logger{
		acc: accChan,
		err: errChan,
	}
	return logger
}

func NewWriterMan() *WriterMan {
	return &WriterMan{
		wg: &sync.WaitGroup{},
	}
}
