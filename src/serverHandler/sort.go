package serverHandler

import (
	"../util"
	"os"
	"sort"
)

type sortableFileInfos struct {
	infos []os.FileInfo
	names [][]byte
}

func newSortableFileInfos(infos []os.FileInfo) sortableFileInfos {
	names := make([][]byte, len(infos))
	for i := range infos {
		names[i] = []byte(infos[i].Name())
	}

	return sortableFileInfos{infos, names}
}

func (sInfos sortableFileInfos) Len() int {
	return len(sInfos.names)
}

func (sInfos sortableFileInfos) Less(i, j int) bool {
	prevIsDir := sInfos.infos[i].IsDir()
	nextIsDir := sInfos.infos[j].IsDir()
	if prevIsDir != nextIsDir {
		return prevIsDir
	}

	return util.CompareNumInStr(sInfos.names[i], sInfos.names[j])
}

func (sInfos sortableFileInfos) Swap(i, j int) {
	sInfos.infos[i], sInfos.infos[j] = sInfos.infos[j], sInfos.infos[i]
	sInfos.names[i], sInfos.names[j] = sInfos.names[j], sInfos.names[i]
}

func sortSubItems(subInfos []os.FileInfo) {
	sortSubInfos := newSortableFileInfos(subInfos)
	sort.Sort(sortSubInfos)
}
