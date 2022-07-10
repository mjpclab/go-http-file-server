package serverHandler

import (
	"../util"
)

type aliasAccurate struct {
	url string
	fs  string
}

func createAliasAccurate(urlPath, fsPath string) aliasAccurate {
	return aliasAccurate{urlPath, fsPath}
}

func (alias aliasAccurate) isMatch(rawReqPath string) bool {
	return util.IsStrEqualAccurate(alias.url, rawReqPath)
}

func (alias aliasAccurate) isSuccessorOf(rawReqPath string) bool {
	return len(alias.url) > len(rawReqPath) && util.HasUrlPrefixDirAccurate(alias.url, rawReqPath)
}

func (alias aliasAccurate) isPredecessorOf(rawReqPath string) bool {
	return len(rawReqPath) > len(alias.url) && util.HasUrlPrefixDirAccurate(rawReqPath, alias.url)
}

func (alias aliasAccurate) namesEqual(a, b string) bool {
	return util.IsStrEqualAccurate(a, b)
}

func (alias aliasAccurate) subPart(rawReqPath string) (subName string, isLastPart, ok bool) {
	if !alias.isSuccessorOf(rawReqPath) {
		return
	}
	subName, isLastPart = getAliasSubPart(alias.url, rawReqPath)
	ok = true
	return
}
