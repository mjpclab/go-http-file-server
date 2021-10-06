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

func createVirtualFileInfo(name string, refItem os.FileInfo, caseSensitive bool) os.FileInfo {
	if refItem != nil {
		if caseSensitive {
			return createRenamedFileInfo(name, refItem)
		} else {
			return createRenamedFileInfoNoCase(name, refItem)
		}
	} else {
		if caseSensitive {
			return createPlaceholderFileInfo(name, true)
		} else {
			return createPlaceholderFileInfoNoCase(name, true)
		}
	}
}

func isVirtual(info os.FileInfo) bool {
	switch info.(type) {
	case placeholderFileInfo, renamedFileInfo, placeholderFileInfoNoCase, renamedFileInfoNoCase:
		return true
	}
	return false
}

func isNameCaseSensitive(info os.FileInfo) bool {
	switch info.(type) {
	case placeholderFileInfoNoCase, renamedFileInfoNoCase:
		return false
	}
	return true
}

func isNameEqualAccurate(a, b string) bool {
	return a == b
}

func isNameEqualNoCase(a, b string) bool {
	return strings.EqualFold(a, b)
}

func getIsNameEqualFunc(info os.FileInfo) func(a, b string) bool {
	if isNameCaseSensitive(info) {
		return isNameEqualAccurate
	} else {
		return isNameEqualNoCase
	}
}

func containsItem(infos []os.FileInfo, name string) bool {
	for i := range infos {
		isNameEqual := getIsNameEqualFunc(infos[i])
		if isNameEqual(infos[i].Name(), name) {
			return true
		}
	}
	return false
}
