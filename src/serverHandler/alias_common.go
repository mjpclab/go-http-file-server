package serverHandler

import "strings"

func getAliasSubPart(aliasPath, rawReqPath string) (subName string, isLastPart bool) {
	subName = aliasPath[len(rawReqPath):]
	if len(subName) > 0 && subName[0] == '/' {
		subName = subName[1:]
	}

	slashIndex := strings.IndexByte(subName, '/')
	if slashIndex > 0 {
		subName = subName[:slashIndex]
	} else {
		isLastPart = true
	}

	return
}
