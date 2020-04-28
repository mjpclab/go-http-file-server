package serverHandler

import (
	"os"
	"time"
)

type dummyFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (info dummyFileInfo) Name() string {
	return info.name
}

func (info dummyFileInfo) Size() int64 {
	return info.size
}

func (info dummyFileInfo) Mode() os.FileMode {
	return 0
}

func (info dummyFileInfo) ModTime() time.Time {
	return info.modTime
}

func (info dummyFileInfo) IsDir() bool {
	return info.isDir
}

func (info dummyFileInfo) Sys() interface{} {
	return nil
}
