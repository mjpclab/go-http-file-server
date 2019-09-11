package util

import (
	"path"
	"strings"
)

var knownSuffixes = [...]string{
	".tar.gz",
	".iso.gz",
	".img.gz",

	".tar.bz2",
	".iso.bz2",
	".img.bz2",

	".tar.xz",
	".iso.xz",
	".img.xz",
}

func SplitFilename(filename string) (prefix, suffix string) {
	if len(filename) == 0 {
		return
	}

	if filename[0] == '.' {
		return filename, ""
	}

	filenameLower := strings.ToLower(filename)
	for _, knownSuffix := range knownSuffixes {
		if len(filenameLower) <= len(knownSuffix) {
			continue
		}

		if filenameLower[len(filenameLower)-len(knownSuffix):] == knownSuffix {
			prefix = filename[:len(filename)-len(knownSuffix)]
			suffix = filename[len(filename)-len(knownSuffix):]
			return prefix, suffix
		}
	}

	suffix = path.Ext(filename)
	prefix = filename[:len(filename)-len(suffix)]
	return prefix, suffix
}
