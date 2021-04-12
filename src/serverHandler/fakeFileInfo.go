package serverHandler

import (
	"os"
	"time"
)

var initTime = time.Now()

type fakeFileInfo struct {
	name  string
	isDir bool
}

func (info *fakeFileInfo) Name() string {
	return info.name
}

func (info *fakeFileInfo) Size() int64 {
	return 0
}

func (info *fakeFileInfo) Mode() os.FileMode {
	return 0
}

func (info *fakeFileInfo) ModTime() time.Time {
	return initTime
}

func (info *fakeFileInfo) IsDir() bool {
	return info.isDir
}

func (info *fakeFileInfo) Sys() interface{} {
	return nil
}

func newFakeFileInfo(name string, isDir bool) *fakeFileInfo {
	return &fakeFileInfo{name, isDir}
}
