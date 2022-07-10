package serverHandler

import "os"

type renamedFileInfoAccurate struct {
	name string
	os.FileInfo
}

func (info renamedFileInfoAccurate) Name() string {
	return info.name
}

func createRenamedFileInfoAccurate(name string, fileInfo os.FileInfo) renamedFileInfoAccurate {
	return renamedFileInfoAccurate{name, fileInfo}
}

type renamedFileInfoNoCase struct {
	renamedFileInfoAccurate
}

func createRenamedFileInfoNoCase(name string, fileInfo os.FileInfo) renamedFileInfoNoCase {
	return renamedFileInfoNoCase{renamedFileInfoAccurate{name, fileInfo}}
}
