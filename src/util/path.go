package util

import (
	"path"
	"path/filepath"
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
	return hasPrefixDir(urlPath, prefix, '/')
}

func HasFsPrefixDir(fsPath, prefix string) bool {
	return hasPrefixDir(fsPath, prefix, filepath.Separator)
}

func hasPrefixDir(absPath, prefix string, separator byte) bool {
	if absPath == prefix {
		return true
	}

	prefixMaxIndex := len(prefix) - 1

	if len(absPath) < len(prefix) {
		if len(absPath) == prefixMaxIndex && prefix[prefixMaxIndex] == separator && absPath == prefix[:prefixMaxIndex] {
			return true
		}
		return false
	}

	if absPath[:len(prefix)] != prefix {
		return false
	}

	if prefix[prefixMaxIndex] == separator {
		return true
	}

	if absPath[len(prefix)] == separator {
		return true
	}

	return false
}
