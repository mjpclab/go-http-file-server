package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
	"strings"
)

type alias struct {
	url string
	dir string
}

func createAlias(urlPath, fsPath string) alias {
	return alias{urlPath, fsPath}
}

func (alias alias) isMatch(rawReqPath string) bool {
	return util.IsPathEqual(alias.url, rawReqPath)
}

func (alias alias) isSuccessorOf(rawReqPath string) bool {
	return len(alias.url) > len(rawReqPath) && util.HasUrlPrefixDir(alias.url, rawReqPath)
}

func (alias alias) isPredecessorOf(rawReqPath string) bool {
	return len(rawReqPath) > len(alias.url) && util.HasUrlPrefixDir(rawReqPath, alias.url)
}

func (alias alias) nextPartOf(rawReqPath string) (subName string, noMore, ok bool) {
	if !alias.isSuccessorOf(rawReqPath) {
		return
	}

	subName = alias.url[len(rawReqPath):]
	if len(subName) > 0 && subName[0] == '/' {
		subName = subName[1:]
	}

	slashIndex := strings.IndexByte(subName, '/')
	if slashIndex > 0 {
		subName = subName[:slashIndex]
	} else {
		noMore = true
	}

	ok = true

	return
}
