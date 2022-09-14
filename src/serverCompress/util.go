package serverCompress

import (
	"mjpclab.dev/ghfs/src/util"
	"strings"
)

var compressibleTypes = []string{
	"application/javascript",
	"application/x-javascript",
	"application/json",
	"application/xhtml+xml",
	"application/xml",
	"image/svg+xml",
}

func isCompressibleType(contentType string) bool {
	if strings.HasPrefix(contentType, "text/") {
		return true
	}

	sepIndex := strings.IndexByte(contentType, ';')
	if sepIndex > 0 {
		contentType = contentType[:sepIndex]
	}
	return util.Contains(compressibleTypes, contentType)
}
