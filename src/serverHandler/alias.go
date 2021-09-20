package serverHandler

import (
	"../util"
	"strings"
)

type alias struct {
	urlPath string
	fsPath  string
}

type aliases []*alias

func NewAlias(urlPath, fsPath string) *alias {
	return &alias{urlPath, fsPath}
}

func NewAliases(capacity int) aliases {
	aliases := make(aliases, 0, capacity)
	return aliases
}

func (aliases aliases) byUrlPath(urlPath string) (alias *alias, ok bool) {
	for _, alias := range aliases {
		if urlPath == alias.urlPath {
			return alias, true
		}
	}
	return nil, false
}

func (alias *alias) isMatch(rawReqPath string) bool {
	return alias.urlPath == rawReqPath
}

func (alias *alias) isSuccessorOf(rawReqPath string) bool {
	return len(alias.urlPath) > len(rawReqPath) && util.HasUrlPrefixDir(alias.urlPath, rawReqPath)
}

func (alias *alias) namesEqual(a, b string) bool {
	return a == b
}

func (alias *alias) getSubPart(rawReqPath string) (subName string, isLastPart, ok bool) {
	if !alias.isSuccessorOf(rawReqPath) {
		return
	}

	subName = alias.urlPath[len(rawReqPath):]
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
