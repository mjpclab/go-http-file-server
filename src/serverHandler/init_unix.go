//go:build !windows
// +build !windows

package serverHandler

import (
	"../util"
)

type alias = aliasAccurate

var createAlias = createAliasAccurate
var isNameEqual = util.IsStrEqualAccurate

var createRenamedFileInfo = createRenamedFileInfoAccurate
var createPlaceholderFileInfo = createPlaceholderFileInfoAccurate
