package serverHandler

import (
	"../util"
)

type aliasNoCase struct {
	url string
	fs  string
}

func createAliasNoCase(urlPath, fsPath string) aliasNoCase {
	return aliasNoCase{urlPath, fsPath}
}

func (alias aliasNoCase) isMatch(rawReqPath string) bool {
	return util.IsStrEqualNoCase(alias.url, rawReqPath)
}

func (alias aliasNoCase) isSuccessorOf(rawReqPath string) bool {
	return len(alias.url) > len(rawReqPath) && util.HasUrlPrefixDirNoCase(alias.url, rawReqPath)
}

func (alias aliasNoCase) isPredecessorOf(rawReqPath string) bool {
	return len(rawReqPath) > len(alias.url) && util.HasUrlPrefixDirNoCase(rawReqPath, alias.url)
}

func (alias aliasNoCase) namesEqual(a, b string) bool {
	return util.IsStrEqualNoCase(a, b)
}

func (alias aliasNoCase) subPart(rawReqPath string) (subName string, isLastPart, ok bool) {
	if !alias.isSuccessorOf(rawReqPath) {
		return
	}
	subName, isLastPart = getAliasSubPart(alias.url, rawReqPath)
	ok = true
	return
}
