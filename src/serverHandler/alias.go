package serverHandler

import (
	"../util"
	"path"
)

type alias struct {
	urlPath string
	fsPath  string
}

type aliases []*alias

func (aliases aliases) byUrlPath(urlPath string) (alias *alias, ok bool) {
	for _, alias := range aliases {
		if urlPath == alias.urlPath {
			return alias, true
		}
	}
	return nil, false
}

func (alias alias) isMatch(rawReqPath string) bool {
	return alias.urlPath == rawReqPath
}

func (alias alias) isChildOf(rawReqPath string) bool {
	return len(alias.urlPath) > len(rawReqPath) && path.Dir(alias.urlPath) == rawReqPath
}

func (alias alias) isSuccessorOf(rawReqPath string) bool {
	return len(alias.urlPath) > len(rawReqPath) && util.HasUrlPrefixDir(alias.urlPath, rawReqPath)
}

func (alias alias) namesEqual(a, b string) bool {
	return a == b
}
