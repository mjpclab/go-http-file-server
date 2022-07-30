package serverHandler

import (
	"os"
	"regexp"
	"testing"
	"time"
)

func TestHandler_FilterItems(t *testing.T) {
	now := time.Now()
	var h *aliasHandler
	var items []os.FileInfo
	var ok bool

	re2 := regexp.MustCompile("2")
	re3 := regexp.MustCompile("3")

	dir1 := dummyFileInfo{"dir1", 0, now, true}
	dir2 := dummyFileInfo{"dir2", 0, now, true}
	dir3 := dummyFileInfo{"dir3", 0, now, true}

	file1 := dummyFileInfo{"file1", 0, now, false}
	file2 := dummyFileInfo{"file2", 0, now, false}
	file3 := dummyFileInfo{"file3", 0, now, false}

	originalItems := []os.FileInfo{dir1, dir2, dir3, file1, file2, file3}

	h = &aliasHandler{shows: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, dir2, file2)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{showDirs: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, dir2, file1, file2, file3)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{showFiles: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, dir1, dir2, dir3, file2)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{hides: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, dir1, dir3, file1, file3)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{hideDirs: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, dir1, dir3, file1, file2, file3)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{hideFiles: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, dir1, dir2, dir3, file1, file3)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{shows: re2, hides: re3}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, dir2, file2)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{shows: re2, hides: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{shows: re2, hideDirs: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, file2)
	if !ok {
		t.Errorf("%+v\n", items)
	}

	h = &aliasHandler{shows: re2, hideFiles: re2}
	items = h.FilterItems(originalItems)
	ok = expectItems(items, dir2)
	if !ok {
		t.Errorf("%+v\n", items)
	}

}
