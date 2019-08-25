package util

import "path"

func CleanUrlPath(urlPath string) string {
	if len(urlPath) == 0 {
		return "/"
	}

	urlPath = path.Clean(urlPath)
	if urlPath[0] != '/' {
		urlPath = "/" + urlPath
	}

	return urlPath
}
