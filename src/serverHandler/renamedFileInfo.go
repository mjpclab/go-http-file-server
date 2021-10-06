package serverHandler

import "os"

type renamedFileInfo struct {
	name string
	os.FileInfo
}

func (info renamedFileInfo) Name() string {
	return info.name
}

func createRenamedFileInfo(name string, fileInfo os.FileInfo) renamedFileInfo {
	return renamedFileInfo{name, fileInfo}
}

type renamedFileInfoNoCase struct {
	renamedFileInfo
}

func createRenamedFileInfoNoCase(name string, fileInfo os.FileInfo) renamedFileInfoNoCase {
	return renamedFileInfoNoCase{renamedFileInfo{name, fileInfo}}
}
