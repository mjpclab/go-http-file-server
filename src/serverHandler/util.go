package serverHandler

import (
	"../acceptHeaders"
	"../util"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func wildcardToRegexp(wildcards []string) (*regexp.Regexp, error) {
	if len(wildcards) == 0 {
		return nil, nil
	}

	normalizedWildcards := make([]string, 0, len(wildcards))
	for _, wildcard := range wildcards {
		if len(wildcard) == 0 {
			continue
		}
		normalizedWildcards = append(normalizedWildcards, util.WildcardToStrRegexp(wildcard))
	}

	if len(normalizedWildcards) == 0 {
		return nil, nil
	}

	exp := strings.Join(normalizedWildcards, "|")
	return regexp.Compile(exp)
}

func needResponseBody(method string) bool {
	return method != http.MethodHead &&
		method != http.MethodOptions &&
		method != http.MethodConnect &&
		method != http.MethodTrace
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
		wr, err = gzip.NewWriterLevel(w, gzip.BestSpeed)
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

func createVirtualFileInfo(name string, refItem os.FileInfo) os.FileInfo {
	if refItem != nil {
		return createRenamedFileInfo(name, refItem)
	} else {
		return createPlaceholderFileInfo(name, true)
	}
}

func isVirtual(info os.FileInfo) bool {
	switch info.(type) {
	case placeholderFileInfo, renamedFileInfo:
		return true
	}
	return false
}

func containsItem(infos []os.FileInfo, name string) bool {
	for i := range infos {
		if util.IsPathEqual(infos[i].Name(), name) {
			return true
		}
	}
	return false
}

func shouldServeAsContent(file *os.File, item os.FileInfo) bool {
	return file != nil && item != nil && !item.IsDir()
}
