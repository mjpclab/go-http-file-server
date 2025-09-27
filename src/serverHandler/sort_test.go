package serverHandler

import (
	"os"
	"testing"
	"time"
)

func expectItems(items []os.FileInfo, expects ...os.FileInfo) bool {
	if len(items) != len(expects) {
		return false
	}

	for i, item := range expects {
		if items[i] != item {
			return false
		}
	}

	return true
}

func TestSort(t *testing.T) {
	now := time.Now()
	var ok bool

	dir1 := dummyFileInfo{"item1", 0, now, true}
	dir2 := dummyFileInfo{"item3", 300, now.Add(time.Minute), true}
	dir3 := dummyFileInfo{"item5", 200, now.Add(time.Minute * 10), true}

	file1 := dummyFileInfo{"item2.zip", 50, now.Add(time.Second), false}
	file2 := dummyFileInfo{"item4.tar", 150, now.Add(time.Minute * 20), false}
	file3 := dummyFileInfo{"item6.zip", 250, now.Add(time.Hour), false}

	originInfos := []os.FileInfo{dir3, file2, dir1, file3, dir2, file1}
	infos := make([]os.FileInfo, len(originInfos))

	copy(infos, originInfos)
	sortInfos(infos, "", "")
	ok = expectItems(infos, dir3, file2, dir1, file3, dir2, file1)
	if !ok {
		t.Error(infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "/n", "")
	ok = expectItems(infos, dir1, dir2, dir3, file1, file2, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "/N", "")
	ok = expectItems(infos, dir3, dir2, dir1, file3, file2, file1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "n/", "")
	ok = expectItems(infos, file1, file2, file3, dir1, dir2, dir3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "N/", "")
	ok = expectItems(infos, file3, file2, file1, dir3, dir2, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "n", "")
	ok = expectItems(infos, dir1, file1, dir2, file2, dir3, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "N", "")
	ok = expectItems(infos, file3, dir3, file2, dir2, file1, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "e", "")
	ok = expectItems(infos, dir1, dir2, dir3, file2, file1, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "E", "")
	ok = expectItems(infos, file3, file1, file2, dir3, dir2, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "s", "")
	ok = expectItems(infos, dir1, file1, file2, dir3, file3, dir2)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "S", "")
	ok = expectItems(infos, dir2, file3, dir3, file2, file1, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "t", "")
	ok = expectItems(infos, dir1, file1, dir2, dir3, file2, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "T", "")
	ok = expectItems(infos, file3, file2, dir3, dir2, file1, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "/", "")
	ok = expectItems(infos, dir3, dir1, dir2, file2, file3, file1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "", "")
	ok = expectItems(infos, originInfos...)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortInfos(infos, "", "/n")
	ok = expectItems(infos, dir1, dir2, dir3, file1, file2, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}
}
