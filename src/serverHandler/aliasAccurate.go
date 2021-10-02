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

func (alias aliasAccurate) urlPath() string {
	return alias.url
}

func (alias aliasAccurate) fsPath() string {
	return alias.fs
}

func (alias aliasAccurate) caseSensitive() bool {
	return true
}

func (alias aliasAccurate) isMatch(rawReqPath string) bool {
	return isNameEqualAccurate(alias.url, rawReqPath)
}

func (alias aliasAccurate) isSuccessorOf(rawReqPath string) bool {
	return len(alias.url) > len(rawReqPath) && util.HasUrlPrefixDir(alias.url, rawReqPath)
}

func (alias aliasAccurate) isPredecessorOf(rawReqPath string) bool {
	return len(rawReqPath) > len(alias.url) && util.HasUrlPrefixDir(rawReqPath, alias.url)
}

func (alias aliasAccurate) namesEqual(a, b string) bool {
	return isNameEqualAccurate(a, b)
}
