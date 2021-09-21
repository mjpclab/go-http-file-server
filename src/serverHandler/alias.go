package serverHandler

import (
	"sort"
	"strings"
)

type alias interface {
	urlPath() string
	fsPath() string
	isMatch(rawReqPath string) bool
	isSuccessorOf(rawReqPath string) bool
	namesEqual(a, b string) bool
}

type aliases []alias

func NewAliases(entriesAccurate, entriesNoCase map[string]string) aliases {
	aliases := make(aliases, 0, len(entriesAccurate)+len(entriesNoCase))
	for urlPath, fsPath := range entriesAccurate {
		aliases = append(aliases, createAliasAccurate(urlPath, fsPath))
	}
	for urlPath, fsPath := range entriesNoCase {
		aliases = append(aliases, createAliasNoCase(urlPath, fsPath))
	}
	sort.Sort(aliases)

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

func (aliases aliases) Len() int {
	return len(aliases)
}

func (aliases aliases) Less(i, j int) bool {
	iLen := len(aliases[i].urlPath())
	jLen := len(aliases[j].urlPath())
	if iLen != jLen {
		// longer is prior
		return iLen > jLen
	}

	_, isIAccurate := aliases[i].(aliasAccurate)
	_, isJAccurate := aliases[j].(aliasAccurate)
	if isIAccurate != isJAccurate {
		// accurate is prior
		return isIAccurate
	}

	return i < j
}

func (aliases aliases) Swap(i, j int) {
	aliases[i], aliases[j] = aliases[j], aliases[i]
}
