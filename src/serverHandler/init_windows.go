//go:build windows
// +build windows

package serverHandler

import (
	"../util"
)

type alias = aliasNoCase

var createAlias = createAliasNoCase
var isNameEqual = util.IsStrEqualNoCase

var createRenamedFileInfo = createRenamedFileInfoNoCase
var createPlaceholderFileInfo = createPlaceholderFileInfoNoCase
