package util

import (
	"path"
	"path/filepath"
	"strings"
)

func CleanUrlPath(urlPath string) string {
	if len(urlPath) == 0 {
		return "/"
	}

	if urlPath[0] != '/' {
		urlPath = "/" + urlPath
	}
	urlPath = path.Clean(urlPath)

	return urlPath
}

func HasUrlPrefixDir(urlPath, prefix string) bool {
	if urlPath == prefix {
		return true
	}

	if prefix[len(prefix)-1] != '/' {
		prefix = prefix + "/"
	}

	return strings.HasPrefix(urlPath, prefix)
}

func HasFsPrefixDir(fsPath, prefix string) bool {
	if fsPath == prefix {
		return true
	}

	if prefix[len(prefix)-1] != filepath.Separator {
		prefix = prefix + string(filepath.Separator)
	}

	return strings.HasPrefix(fsPath, prefix)
}
