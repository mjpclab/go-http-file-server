package serverHandler

import (
	"../util"
	"strings"
)

type aliasNoCase struct {
	url string
	fs  string
}

func createAliasNoCase(urlPath, fsPath string) aliasNoCase {
	return aliasNoCase{urlPath, fsPath}
}

func (alias aliasNoCase) urlPath() string {
	return alias.url
}

func (alias aliasNoCase) fsPath() string {
	return alias.fs
}

func (alias aliasNoCase) isMatch(rawReqPath string) bool {
	return strings.EqualFold(alias.url, rawReqPath)
}

func (alias aliasNoCase) isSuccessorOf(rawReqPath string) bool {
	return len(alias.url) > len(rawReqPath) && util.HasUrlPrefixDirNoCase(alias.url, rawReqPath)
}

func (alias aliasNoCase) namesEqual(a, b string) bool {
	return strings.EqualFold(a, b)
}
