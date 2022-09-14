package serverCompress

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"mjpclab.dev/ghfs/src/acceptHeaders"
	"net/http"
)

const contentEncGzip = "gzip"
const contentEncDeflate = "deflate"

var encodings = []string{contentEncGzip, contentEncDeflate}

func GetWriter(w http.ResponseWriter, r *http.Request) (wc io.WriteCloser, ok bool) {
	header := w.Header()
	if len(header.Get("Content-Encoding")) > 0 {
		return nil, false
	}
	if !isCompressibleType(header.Get("Content-Type")) {
		return nil, false
	}

	accepts := acceptHeaders.ParseAccepts(r.Header.Get("Accept-Encoding"))
	_, encoding, hasSupportedEncoding := accepts.GetPreferredValue(encodings)
	if !hasSupportedEncoding {
		return nil, false
	}

	var err error
	switch encoding {
	case contentEncGzip:
		wc, err = gzip.NewWriterLevel(w, gzip.BestSpeed)
	case contentEncDeflate:
		wc, err = flate.NewWriter(w, flate.BestSpeed)
	default:
		return nil, false
	}

	if err != nil {
		return nil, false
	}

	header.Del("Content-Length")
	header.Set("Content-Encoding", encoding)
	return wc, true
}
