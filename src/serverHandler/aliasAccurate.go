package serverHandler

import (
	"../util"
	"strings"
)

type aliasAccurate struct {
	url string
	fs  string
}

func CreateAliasAccurate(urlPath, fsPath string) aliasAccurate {
	return aliasAccurate{urlPath, fsPath}
}

func (alias aliasAccurate) urlPath() string {
	return alias.url
}

func (alias aliasAccurate) fsPath() string {
	return alias.fs
}

func (alias aliasAccurate) isMatch(rawReqPath string) bool {
	return alias.url == rawReqPath
}

func (alias aliasAccurate) isSuccessorOf(rawReqPath string) bool {
	return len(alias.url) > len(rawReqPath) && util.HasUrlPrefixDir(alias.url, rawReqPath)
}

func (alias aliasAccurate) namesEqual(a, b string) bool {
	return a == b
}

func (alias aliasAccurate) getSubPart(rawReqPath string) (subName string, isLastPart, ok bool) {
	if !alias.isSuccessorOf(rawReqPath) {
		return
	}

	subName = alias.url[len(rawReqPath):]
	if len(subName) > 0 && subName[0] == '/' {
		subName = subName[1:]
	}

	slashIndex := strings.Index(subName, "/")
	if slashIndex > 0 {
		subName = subName[:slashIndex]
	} else {
		isLastPart = true
	}

	ok = true

	return
}
