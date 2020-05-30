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

// infosNames

type infosNames struct {
	infos      []os.FileInfo
	names      [][]byte
	compareDir fnCompareDir
}

func (xInfos infosNames) Len() int {
	return len(xInfos.names)
}

func (xInfos infosNames) Swap(i, j int) {
	xInfos.infos[i], xInfos.infos[j] = xInfos.infos[j], xInfos.infos[i]
	xInfos.names[i], xInfos.names[j] = xInfos.names[j], xInfos.names[i]
}

func (xInfos infosNames) LessDir(i, j int) (less, ok bool) {
	return xInfos.compareDir(xInfos.infos[i], xInfos.infos[j])
}

func (xInfos infosNames) LessFilename(i, j int) (less, ok bool) {
	return util.CompareNumInFilename(xInfos.names[i], xInfos.names[j])
}

func newInfosNames(infos []os.FileInfo, compareDir fnCompareDir) infosNames {
	names := make([][]byte, len(infos))
	for i := range infos {
		names[i] = []byte(infos[i].Name())
	}

	return infosNames{infos, names, compareDir}
}

// sort name asc

type infosNamesAsc struct {
	infosNames
}

func (xInfos infosNamesAsc) Less(i, j int) bool {
	less, ok := xInfos.LessDir(i, j)
	if ok {
		return less
	}

	less, ok = xInfos.LessFilename(i, j)
	if ok {
		return less
	}

	return i < j
}

func sortInfoNamesAsc(infos []os.FileInfo, compareDir fnCompareDir) {
	nameCachedInfos := infosNamesAsc{newInfosNames(infos, compareDir)}
	sort.Sort(nameCachedInfos)
}

// sort name desc

type infosNamesDesc struct {
	infosNames
}

func (xInfos infosNamesDesc) Less(i, j int) bool {
	less, ok := xInfos.LessDir(i, j)
	if ok {
		return less
	}

	less, ok = xInfos.LessFilename(j, i)
	if ok {
		return less
	}

	return j < i
}

func sortInfoNamesDesc(infos []os.FileInfo, compareDir fnCompareDir) {
	nameCachedInfos := infosNamesDesc{newInfosNames(infos, compareDir)}
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
