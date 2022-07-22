package serverHandler

import (
	"os"
	"testing"
)

func TestGetCleanFilePath(t *testing.T) {
	var cleanPath string
	var ok bool

	cleanPath, ok = getCleanFilePath("file1")
	if cleanPath != "file1" || !ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanFilePath("./file2")
	if cleanPath != "file2" || !ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanFilePath("./dir/file3")
	if ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanFilePath("dir/file4")
	if ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanFilePath("../file5")
	if ok {
		t.Error(cleanPath, ok)
	}
}

func TestGetCleanDirFilePath(t *testing.T) {
	var cleanPath string
	var ok bool

	cleanPath, ok = getCleanDirFilePath("file1")
	if cleanPath != "file1" || !ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanDirFilePath("./file2")
	if cleanPath != "file2" || !ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanDirFilePath("./dir/file3")
	if cleanPath != "dir/file3" || !ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanDirFilePath("dir/file4")
	if cleanPath != "dir/file4" || !ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanDirFilePath("dir1/../dir2/file5")
	if cleanPath != "dir2/file5" || !ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanDirFilePath("dir1/../../dir2/file6")
	if ok {
		t.Error(cleanPath, ok)
	}

	cleanPath, ok = getCleanDirFilePath("../file5")
	if ok {
		t.Error(cleanPath, ok)
	}
}

func TestIsVirtual(t *testing.T) {
	var info os.FileInfo

	info = createPlaceholderFileInfo("foo", true)
	if !isVirtual(info) {
		t.Error()
	}

	baseInfo := dummyFileInfo{name: "foo"}
	if isVirtual(baseInfo) {
		t.Error()
	}

	info = createRenamedFileInfo("bar", baseInfo)
	if !isVirtual(info) {
		t.Error()
	}
}
