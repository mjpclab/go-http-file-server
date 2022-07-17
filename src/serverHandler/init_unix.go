//go:build !windows
// +build !windows

package serverHandler

type alias = aliasAccurate

var createAlias = createAliasAccurate

var createRenamedFileInfo = createRenamedFileInfoAccurate
var createPlaceholderFileInfo = createPlaceholderFileInfoAccurate
