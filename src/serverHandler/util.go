package serverHandler

import (
	"net/http"
	"os"
	"path"
	"strings"
)

func needResponseBody(method string) bool {
	return method != http.MethodHead &&
		method != http.MethodOptions &&
		method != http.MethodConnect &&
		method != http.MethodTrace
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
	ok = filePath != "." && filePath != ".." &&
		strings.IndexByte(filePath, '/') < 0 &&
		(os.PathSeparator == '/' || strings.IndexByte(filePath, os.PathSeparator) < 0)

	return
}

func getCleanDirFilePath(requestPath string) (filePath string, ok bool) {
	filePath = path.Clean(strings.Replace(requestPath, "\\", "/", -1))
	ok = filePath[0] != '/' && filePath != "." && filePath != ".." && !strings.HasPrefix(filePath, "../")

	return
}
