package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
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

func getRedirectCode(r *http.Request) int {
	if r.Method == http.MethodPost {
		return http.StatusTemporaryRedirect
	} else {
		return http.StatusMovedPermanently
	}
}

func NeedResponseBody(method string) bool {
	return method != http.MethodHead &&
		method != http.MethodOptions &&
		method != http.MethodConnect &&
		method != http.MethodTrace
}

func lacksHeader(header http.Header, key string) bool {
	return len(header.Get(key)) == 0
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
