package serverHandler

import (
	"os"
	"time"
)

var initTime = time.Now()

type placeholderFileInfo struct {
	name  string
	isDir bool
}

func (info placeholderFileInfo) Name() string {
	return info.name
}

func (info placeholderFileInfo) Size() int64 {
	return 0
}

func (info placeholderFileInfo) Mode() os.FileMode {
	return 0
}

func (info placeholderFileInfo) ModTime() time.Time {
	return initTime
}

func (info placeholderFileInfo) IsDir() bool {
	return info.isDir
}

func (info placeholderFileInfo) Sys() interface{} {
	return nil
}

func createPlaceholderFileInfo(name string, isDir bool) placeholderFileInfo {
	return placeholderFileInfo{name, isDir}
}
