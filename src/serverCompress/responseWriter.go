package serverCompress

import (
	"io"
	"net/http"
)

const minLengthDigits = 4

// responseWriter implements http.ResponseWriter
type responseWriter struct {
	rw     http.ResponseWriter
	r      *http.Request
	writer io.Writer
	closer io.Closer
}

func (rw *responseWriter) Header() http.Header {
	return rw.rw.Header()
}

func (rw *responseWriter) Write(bs []byte) (int, error) {
	if rw.writer == nil {
		rw.init(0)
	}
	return rw.writer.Write(bs)
}

func (rw *responseWriter) WriteHeader(status int) {
	if rw.writer == nil {
		rw.init(status)
	}
	rw.rw.WriteHeader(status)
}

func (rw *responseWriter) init(status int) {
	lengthDigits := len(rw.rw.Header().Get("Content-Length"))
	if lengthDigits >= minLengthDigits || (lengthDigits == 0 && status/100 != 3) {
		wc, ok := GetWriter(rw.rw, rw.r)
		if ok {
			rw.writer = wc
			rw.closer = wc
			return
		}
	}

	rw.writer = rw.rw
}

func (rw *responseWriter) Close() {
	if rw.closer != nil {
		rw.closer.Close()
	}
}

func NewResponseWriter(w http.ResponseWriter, r *http.Request) *responseWriter {
	return &responseWriter{
		rw: w,
		r:  r,
	}
}
