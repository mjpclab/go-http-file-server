package util

import (
	"path"
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

func hasPrefixDirAccurate(absPath, prefix string, separator byte) bool {
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

func hasPrefixDirNoCase(absPath, prefix string, separator byte) bool {
	if len(absPath) == len(prefix) {
		return strings.EqualFold(absPath, prefix)
	}

	prefixMaxIndex := len(prefix) - 1

	if len(absPath) < len(prefix) {
		if len(absPath) == prefixMaxIndex && prefix[prefixMaxIndex] == separator && strings.EqualFold(absPath, prefix[:prefixMaxIndex]) {
			return true
		}
		return false
	}

	if !strings.EqualFold(absPath[:len(prefix)], prefix) {
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

func NormalizeUrlPath(input string) (string, error) {
	result := CleanUrlPath(input)
	return result, nil
}
