package serverHandler

import (
	"../acceptHeaders"
	"../shimgo"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func needResponseBody(method string) bool {
	return method != shimgo.Net_Http_MethodHead &&
		method != shimgo.Net_Http_MethodOptions &&
		method != shimgo.Net_Http_MethodConnect &&
		method != shimgo.Net_Http_MethodTrace
}

func containsItem(infos []os.FileInfo, name string) bool {
	for i := range infos {
		if infos[i].Name() == name {
			return true
		}
	}
	return false
}

func getCleanFilePath(requestPath string) (filePath string, ok bool) {
	filePath = path.Clean(requestPath)
	ok = filePath == path.Base(filePath)

	return
}

func getCleanDirFilePath(requestPath string) (filePath string, ok bool) {
	filePath = path.Clean(strings.Replace(requestPath, "\\", "/", -1))
	ok = filePath[0] != '/' && filePath != "." && filePath != ".." && !strings.HasPrefix(filePath, "../")

	return
}

const contentEncGzip = "gzip"
const contentEncDeflate = "deflate"

var encodings = []string{contentEncGzip, contentEncDeflate}

func getCompressWriter(w http.ResponseWriter, r *http.Request) (wr io.WriteCloser, encoding string, ok bool) {
	accepts := acceptHeaders.ParseAccepts(r.Header.Get("Accept-Encoding"))
	_, encoding, ok = accepts.GetPreferredValue(encodings)
	if !ok {
		return nil, "", false
	}

	var err error
	switch encoding {
	case contentEncGzip:
		wr, err = gzip.NewWriterLevel(w, flate.BestSpeed)
	case contentEncDeflate:
		wr, err = flate.NewWriter(w, flate.BestSpeed)
	default:
		return nil, "", false
	}

	if err != nil {
		return nil, "", false
	}
	return wr, encoding, true
}

func isVirtual(info os.FileInfo) bool {
	if _, isRenamedInfo := info.(*renamedFileInfo); isRenamedInfo {
		return true
	}
	if _, isFakeInfo := info.(*fakeFileInfo); isFakeInfo {
		return true
	}
	return false
}
