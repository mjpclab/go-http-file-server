package serverHandler

import (
	"os"
	"time"
)

var initTime = time.Now()

type placeholderFileInfoAccurate struct {
	name  string
	isDir bool
}

func (info placeholderFileInfoAccurate) Name() string {
	return info.name
}

func (info placeholderFileInfoAccurate) Size() int64 {
	return 0
}

func (info placeholderFileInfoAccurate) Mode() os.FileMode {
	return 0
}

func (info placeholderFileInfoAccurate) ModTime() time.Time {
	return initTime
}

func (info placeholderFileInfoAccurate) IsDir() bool {
	return info.isDir
}

func (info placeholderFileInfoAccurate) Sys() interface{} {
	return nil
}

func createPlaceholderFileInfoAccurate(name string, isDir bool) placeholderFileInfoAccurate {
	return placeholderFileInfoAccurate{name, isDir}
}

type placeholderFileInfoNoCase struct {
	placeholderFileInfoAccurate
}

func createPlaceholderFileInfoNoCase(name string, isDir bool) placeholderFileInfoNoCase {
	return placeholderFileInfoNoCase{placeholderFileInfoAccurate{name, isDir}}
}
