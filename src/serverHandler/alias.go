package serverHandler

import "strings"

type alias interface {
	urlPath() string
	fsPath() string
	isMatch(rawReqPath string) bool
	isSuccessorOf(rawReqPath string) bool
	namesEqual(a, b string) bool
}

type aliases []alias

func NewAliases(capacity int) aliases {
	aliases := make(aliases, 0, capacity)
	return aliases
}

func (aliases aliases) byUrlPath(urlPath string) (alias alias, ok bool) {
	for _, alias := range aliases {
		if alias.isMatch(urlPath) {
			return alias, true
		}
	}
	return nil, false
}

func getAliasSubPart(alias alias, rawReqPath string) (subName string, isLastPart, ok bool) {
	if !alias.isSuccessorOf(rawReqPath) {
		return
	}

	subName = alias.urlPath()[len(rawReqPath):]
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
