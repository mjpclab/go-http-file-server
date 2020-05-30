package serverHandler

import (
	"../util"
	"os"
	"sort"
	"strings"
)

// compare dir func

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

// sort name asc

type sortNameAscInfos struct {
	infos      []os.FileInfo
	names      [][]byte
	compareDir fnCompareDir
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

func newSortNameAscInfos(infos []os.FileInfo, compareDir fnCompareDir) sortNameAscInfos {
	names := make([][]byte, len(infos))
	for i := range infos {
		names[i] = []byte(infos[i].Name())
	}

	return sortNameAscInfos{infos, names, compareDir}
}

func sortInfoNamesAsc(infos []os.FileInfo, compareDir fnCompareDir) {
	nameCachedInfos := newSortNameAscInfos(infos, compareDir)
	sort.Sort(nameCachedInfos)
}

// sort name desc

type sortNameDescInfos struct {
	sortNameAscInfos
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

func newSortNameDescInfos(infos []os.FileInfo, compareDir fnCompareDir) sortNameDescInfos {
	return sortNameDescInfos{newSortNameAscInfos(infos, compareDir)}
}

func sortInfoNamesDesc(infos []os.FileInfo, compareDir fnCompareDir) {
	nameCachedInfos := newSortNameDescInfos(infos, compareDir)
	sort.Sort(nameCachedInfos)
}

// sort size asc

func sortInfoSizesAsc(infos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(infos, func(i, j int) bool {
		less, ok := compareDir(infos[i], infos[j])
		if ok {
			return less
		}

		if infos[i].Size() != infos[j].Size() {
			return infos[i].Size() < infos[j].Size()
		}

		cmpResult := strings.Compare(infos[i].Name(), infos[j].Name())
		if cmpResult != 0 {
			return cmpResult < 0
		}

		return i < j
	})
}

// sort size desc

func sortInfoSizesDesc(infos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(infos, func(i, j int) bool {
		less, ok := compareDir(infos[i], infos[j])
		if ok {
			return less
		}

		if infos[j].Size() != infos[i].Size() {
			return infos[j].Size() < infos[i].Size()
		}

		cmpResult := strings.Compare(infos[j].Name(), infos[i].Name())
		if cmpResult != 0 {
			return cmpResult < 0
		}

		return j < i
	})
}

// sort time asc

func sortInfoTimesAsc(infos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(infos, func(i, j int) bool {
		less, ok := compareDir(infos[i], infos[j])
		if ok {
			return less
		}

		if !infos[i].ModTime().Equal(infos[j].ModTime()) {
			return infos[i].ModTime().Before(infos[j].ModTime())
		}

		cmpResult := strings.Compare(infos[i].Name(), infos[j].Name())
		if cmpResult != 0 {
			return cmpResult < 0
		}

		return i < j
	})
}

// sort time desc

func sortInfoTimesDesc(infos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(infos, func(i, j int) bool {
		less, ok := compareDir(infos[i], infos[j])
		if ok {
			return less
		}

		if !infos[j].ModTime().Equal(infos[i].ModTime()) {
			return infos[j].ModTime().Before(infos[i].ModTime())
		}

		cmpResult := strings.Compare(infos[j].Name(), infos[i].Name())
		if cmpResult != 0 {
			return cmpResult < 0
		}

		return j < i
	})
}

// sort original

func sortInfoOriginal(infos []os.FileInfo, compareDir fnCompareDir) {
	sort.Slice(infos, func(i, j int) bool {
		less, ok := compareDir(infos[i], infos[j])
		if ok {
			return less
		}

		return i < j
	})
}

// sort

func sortInfos(infos []os.FileInfo, rawQuery string, defaultSortBy string) (rawSortBy *string, sortInfo SortState) {
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
		sortInfoNamesAsc(infos, compareDir)
	case 'N':
		sortInfoNamesDesc(infos, compareDir)
	case 's':
		sortInfoSizesAsc(infos, compareDir)
	case 'S':
		sortInfoSizesDesc(infos, compareDir)
	case 't':
		sortInfoTimesAsc(infos, compareDir)
	case 'T':
		sortInfoTimesDesc(infos, compareDir)
	default:
		if dirSort != dirSortMixed {
			sortInfoOriginal(infos, compareDir)
		}
	}

	return rawSortBy, SortState{dirSort, sortKey}
}
