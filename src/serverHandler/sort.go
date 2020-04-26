package serverHandler

import (
	"../util"
	"os"
	"sort"
	"strings"
)

type fnCompareDir func(i, j os.FileInfo) (less, ok bool)

var cmpDirFirst fnCompareDir = func(prev, next os.FileInfo) (less, ok bool) {
	prevIsDir := prev.IsDir()
	nextIsDir := next.IsDir()
	if prevIsDir != nextIsDir {
		return prevIsDir, true
	}
	return true, false
}

var cmpDirLast fnCompareDir = func(prev, next os.FileInfo) (less, ok bool) {
	prevIsDir := prev.IsDir()
	nextIsDir := next.IsDir()
	if prevIsDir != nextIsDir {
		return !prevIsDir, true
	}
	return true, false
}

var cmpDirMixed fnCompareDir = func(prev, next os.FileInfo) (less, ok bool) {
	return true, false
}

type sortNameAscInfos struct {
	infos      []os.FileInfo
	names      [][]byte
	compareDir fnCompareDir
}

func newSortNameAscInfos(infos []os.FileInfo, compareDir fnCompareDir) sortNameAscInfos {
	names := make([][]byte, len(infos))
	for i := range infos {
		names[i] = []byte(infos[i].Name())
	}

	return sortNameAscInfos{infos, names, compareDir}
}

func (sInfos sortNameAscInfos) Len() int {
	return len(sInfos.names)
}

func (sInfos sortNameAscInfos) Less(i, j int) bool {
	less, ok := sInfos.compareDir(sInfos.infos[i], sInfos.infos[j])
	if ok {
		return less
	}

	less, ok = util.CompareNumInFilename(sInfos.names[i], sInfos.names[j])
	if ok {
		return less
	}

	return i < j
}

func (sInfos sortNameAscInfos) Swap(i, j int) {
	sInfos.infos[i], sInfos.infos[j] = sInfos.infos[j], sInfos.infos[i]
	sInfos.names[i], sInfos.names[j] = sInfos.names[j], sInfos.names[i]
}

type sortNameDescInfos struct {
	sortNameAscInfos
}

func newSortNameDescInfos(infos []os.FileInfo, compareDir fnCompareDir) sortNameDescInfos {
	return sortNameDescInfos{newSortNameAscInfos(infos, compareDir)}
}

func (sInfos sortNameDescInfos) Less(i, j int) bool {
	less, ok := sInfos.compareDir(sInfos.infos[i], sInfos.infos[j])
	if ok {
		return less
	}

	less, ok = util.CompareNumInFilename(sInfos.names[j], sInfos.names[i])
	if ok {
		return less
	}

	return j < i
}

func sortSubItemNamesAsc(subInfos []os.FileInfo, compareDir fnCompareDir) {
	nameCachedSubInfos := newSortNameAscInfos(subInfos, compareDir)
	sort.Sort(nameCachedSubInfos)
}

func sortSubItemNamesDesc(subInfos []os.FileInfo, compareDir fnCompareDir) {
	nameCachedSubInfos := newSortNameDescInfos(subInfos, compareDir)
	sort.Sort(nameCachedSubInfos)
}

func sortSubItemSizesAsc(subInfos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(subInfos, func(i, j int) bool {
		less, ok := compareDir(subInfos[i], subInfos[j])
		if ok {
			return less
		}

		if subInfos[i].Size() != subInfos[j].Size() {
			return subInfos[i].Size() < subInfos[j].Size()
		}

		cmpResult := strings.Compare(subInfos[i].Name(), subInfos[j].Name())
		if cmpResult != 0 {
			return cmpResult < 0
		}

		return i < j
	})
}

func sortSubItemSizesDesc(subInfos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(subInfos, func(i, j int) bool {
		less, ok := compareDir(subInfos[i], subInfos[j])
		if ok {
			return less
		}

		if subInfos[j].Size() != subInfos[i].Size() {
			return subInfos[j].Size() < subInfos[i].Size()
		}

		cmpResult := strings.Compare(subInfos[j].Name(), subInfos[i].Name())
		if cmpResult != 0 {
			return cmpResult < 0
		}

		return j < i
	})
}

func sortSubItemTimesAsc(subInfos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(subInfos, func(i, j int) bool {
		less, ok := compareDir(subInfos[i], subInfos[j])
		if ok {
			return less
		}

		if !subInfos[i].ModTime().Equal(subInfos[j].ModTime()) {
			return subInfos[i].ModTime().Before(subInfos[j].ModTime())
		}

		cmpResult := strings.Compare(subInfos[i].Name(), subInfos[j].Name())
		if cmpResult != 0 {
			return cmpResult < 0
		}

		return i < j
	})
}

func sortSubItemTimesDesc(subInfos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(subInfos, func(i, j int) bool {
		less, ok := compareDir(subInfos[i], subInfos[j])
		if ok {
			return less
		}

		if !subInfos[j].ModTime().Equal(subInfos[i].ModTime()) {
			return subInfos[j].ModTime().Before(subInfos[i].ModTime())
		}

		cmpResult := strings.Compare(subInfos[j].Name(), subInfos[i].Name())
		if cmpResult != 0 {
			return cmpResult < 0
		}

		return j < i
	})
}

func sortSubItemOriginal(subInfos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(subInfos, func(i, j int) bool {
		less, ok := compareDir(subInfos[i], subInfos[j])
		if ok {
			return less
		}

		return i < j
	})
}
func sortSubItems(subInfos []os.FileInfo, rawQuery string, defaultSortBy string) (rawSortBy *string, sortInfo SortState) {
	const sortPrefix = "sort="
	var sortBy string

	// extract sort string
	iSortBy := strings.Index(rawQuery, sortPrefix)
	if iSortBy < 0 {
		sortBy = defaultSortBy
	} else {
		if len(rawQuery) > iSortBy+len(sortPrefix) {
			sortBy = rawQuery[iSortBy+len(sortPrefix):]
			iAmp := strings.IndexByte(sortBy, '&')
			if iAmp >= 0 {
				sortBy = sortBy[:iAmp]
			}
		}
		rawSortBy = &sortBy
	}

	if len(sortBy) == 0 {
		return
	}

	// prepare sort info
	var dirSort dirSort
	var sortKey byte
	var compareDir fnCompareDir
	switch {
	case sortBy[0] == '/':
		dirSort = dirSortFirst
		if len(sortBy) > 1 {
			sortKey = sortBy[1]
		}
		compareDir = cmpDirFirst
	case sortBy[len(sortBy)-1] == '/':
		dirSort = dirSortLast
		sortKey = sortBy[0]
		compareDir = cmpDirLast
	default:
		dirSort = dirSortMixed
		sortKey = sortBy[0]
		compareDir = cmpDirMixed
	}

	// do sort
	switch sortKey {
	case 'n':
		sortSubItemNamesAsc(subInfos, compareDir)
	case 'N':
		sortSubItemNamesDesc(subInfos, compareDir)
	case 's':
		sortSubItemSizesAsc(subInfos, compareDir)
	case 'S':
		sortSubItemSizesDesc(subInfos, compareDir)
	case 't':
		sortSubItemTimesAsc(subInfos, compareDir)
	case 'T':
		sortSubItemTimesDesc(subInfos, compareDir)
	default:
		if dirSort != dirSortMixed {
			sortSubItemOriginal(subInfos, compareDir)
		}
	}

	return rawSortBy, SortState{dirSort, sortKey}
}
