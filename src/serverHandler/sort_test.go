package serverHandler

import (
	"os"
	"testing"
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

func expectInfos(infos []os.FileInfo, expects ...os.FileInfo) bool {
	if len(infos) != len(expects) {
		return false
	}

	for i, info := range expects {
		if infos[i] != info {
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

	file1 := dummyFileInfo{"item2", 50, now.Add(time.Second), false}
	file2 := dummyFileInfo{"item4", 150, now.Add(time.Minute * 20), false}
	file3 := dummyFileInfo{"item6", 250, now.Add(time.Hour), false}

	originInfos := []os.FileInfo{dir3, file2, dir1, file3, dir2, file1}
	infos := make([]os.FileInfo, len(originInfos))

	copy(infos, originInfos)
	sortSubItems(infos, "", "")
	ok = expectInfos(infos, dir3, file2, dir1, file3, dir2, file1)
	if !ok {
		t.Error(infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=/n", "")
	ok = expectInfos(infos, dir1, dir2, dir3, file1, file2, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=/N", "")
	ok = expectInfos(infos, dir3, dir2, dir1, file3, file2, file1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=n/", "")
	ok = expectInfos(infos, file1, file2, file3, dir1, dir2, dir3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=N/", "")
	ok = expectInfos(infos, file3, file2, file1, dir3, dir2, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=n", "")
	ok = expectInfos(infos, dir1, file1, dir2, file2, dir3, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=N", "")
	ok = expectInfos(infos, file3, dir3, file2, dir2, file1, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=s", "")
	ok = expectInfos(infos, dir1, file1, file2, dir3, file3, dir2)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=S", "")
	ok = expectInfos(infos, dir2, file3, dir3, file2, file1, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=t", "")
	ok = expectInfos(infos, dir1, file1, dir2, dir3, file2, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=T", "")
	ok = expectInfos(infos, file3, file2, dir3, dir2, file1, dir1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "?sort=/", "")
	ok = expectInfos(infos, dir3, dir1, dir2, file2, file3, file1)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "", "")
	ok = expectInfos(infos, originInfos...)
	if !ok {
		t.Errorf("%+v\n", infos)
	}

	copy(infos, originInfos)
	sortSubItems(infos, "", "/n")
	ok = expectInfos(infos, dir1, dir2, dir3, file1, file2, file3)
	if !ok {
		t.Errorf("%+v\n", infos)
	}
}
