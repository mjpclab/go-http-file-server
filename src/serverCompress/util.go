package serverCompress

import (
	"mjpclab.dev/ghfs/src/util"
	"strings"
)

var compressibleTypes = []string{
	"application/javascript",
	"application/x-javascript",
	"application/json",
	"application/xml",
}

func isCompressibleType(contentType string) bool {
	if strings.HasPrefix(contentType, "text/") {
		return true
	}

	sepIndex := strings.IndexByte(contentType, ';')
	if sepIndex > 0 {
		contentType = contentType[:sepIndex]
	}

	// "image/svg+xml", "application/xhtml+xml", ...
	if strings.HasSuffix(contentType, "+xml") {
		return true
	}

	return util.Contains(compressibleTypes, contentType)
}
