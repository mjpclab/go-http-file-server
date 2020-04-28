package serverHandler

import (
	"os"
	"regexp"
	"testing"
	"time"
)

func TestHandler_FilterItems(t *testing.T) {
	now := time.Now()
	var h *handler
	var infos []os.FileInfo
	var ok bool

	re2 := regexp.MustCompile("2")
	re3 := regexp.MustCompile("3")

	dir1 := dummyFileInfo{"dir1", 0, now, true}
	dir2 := dummyFileInfo{"dir2", 0, now, true}
	dir3 := dummyFileInfo{"dir3", 0, now, true}

	file1 := dummyFileInfo{"file1", 0, now, false}
	file2 := dummyFileInfo{"file2", 0, now, false}
	file3 := dummyFileInfo{"file3", 0, now, false}

	originInfos := []os.FileInfo{dir1, dir2, dir3, file1, file2, file3}

	h = &handler{shows: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, dir2, file2)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{showDirs: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, dir2, file1, file2, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{showFiles: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, dir1, dir2, dir3, file2)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{hides: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, dir1, dir3, file1, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{hideDirs: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, dir1, dir3, file1, file2, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{hideFiles: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, dir1, dir2, dir3, file1, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{shows: re2, hides: re3}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, dir2, file2)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{shows: re2, hides: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{shows: re2, hideDirs: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, file2)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	h = &handler{shows: re2, hideFiles: re2}
	infos = h.FilterItems(originInfos)
	ok = expectInfos(infos, dir2)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

}
