package serverHandler

import (
	"bytes"
	"mjpclab.dev/ghfs/src/util"
	"os"
	"sort"
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
		return nextIsDir, true
	}
	return true, false
}

var cmpDirMixed fnCompareDir = func(prev, next os.FileInfo) (less, ok bool) {
	return true, false
}

// infos
type infos struct {
	items      []os.FileInfo
	compareDir fnCompareDir
}

func (infos infos) Len() int {
	return len(infos.items)
}

func (infos infos) Swap(i, j int) {
	infos.items[i], infos.items[j] = infos.items[j], infos.items[i]
}

func newInfos(items []os.FileInfo, compareDir fnCompareDir) infos {
	return infos{items, compareDir}
}

// infosNames

type infosNames struct {
	items      []os.FileInfo
	names      [][]byte
	compareDir fnCompareDir
}

func (xInfos infosNames) Len() int {
	return len(xInfos.names)
}

func (xInfos infosNames) Swap(i, j int) {
	xInfos.items[i], xInfos.items[j] = xInfos.items[j], xInfos.items[i]
	xInfos.names[i], xInfos.names[j] = xInfos.names[j], xInfos.names[i]
}

func (xInfos infosNames) LessDir(i, j int) (less, ok bool) {
	return xInfos.compareDir(xInfos.items[i], xInfos.items[j])
}

func (xInfos infosNames) LessFilename(i, j int) (less, ok bool) {
	return util.CompareNumInFilename(xInfos.names[i], xInfos.names[j])
}

func (xInfos infosNames) LessFileType(i, j int) (less, ok bool) {
	bufferI := xInfos.names[i]
	bufferJ := xInfos.names[j]
	for {
		dotIndexI := bytes.LastIndexByte(bufferI, '.')
		dotIndexJ := bytes.LastIndexByte(bufferJ, '.')
		if dotIndexI < 0 && dotIndexJ < 0 {
			break
		}
		if dotIndexI < 0 {
			return true, true
		}
		if dotIndexJ < 0 {
			return false, true
		}

		typeI := bufferI[dotIndexI+1:]
		typeJ := bufferJ[dotIndexJ+1:]
		less, ok = util.CompareNumInFilename(typeI, typeJ)
		if ok {
			return less, ok
		}
		bufferI = bufferI[:dotIndexI]
		bufferJ = bufferJ[:dotIndexJ]
	}

	return util.CompareNumInFilename(bufferI, bufferJ)
}

func newInfosNames(items []os.FileInfo, compareDir fnCompareDir) infosNames {
	names := make([][]byte, len(items))
	for i := range items {
		names[i] = []byte(items[i].Name())
	}

	return infosNames{items, names, compareDir}
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

func sortInfoNamesAsc(items []os.FileInfo, compareDir fnCompareDir) {
	nameCachedInfos := infosNamesAsc{newInfosNames(items, compareDir)}
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

func sortInfoNamesDesc(items []os.FileInfo, compareDir fnCompareDir) {
	nameCachedInfos := infosNamesDesc{newInfosNames(items, compareDir)}
	sort.Sort(nameCachedInfos)
}

// sort type asc

type infosTypesAsc struct {
	infosNames
}

func (xInfos infosTypesAsc) Less(i, j int) bool {
	less, ok := xInfos.LessDir(i, j)
	if ok {
		return less
	}

	less, ok = xInfos.LessFileType(i, j)
	if ok {
		return less
	}

	return i < j
}

func sortInfoTypesAsc(items []os.FileInfo, compareDir fnCompareDir) {
	nameCachedInfos := infosTypesAsc{newInfosNames(items, compareDir)}
	sort.Sort(nameCachedInfos)
}

// sort type desc

type infosTypesDesc struct {
	infosNames
}

func (xInfos infosTypesDesc) Less(i, j int) bool {
	less, ok := xInfos.LessDir(i, j)
	if ok {
		return less
	}

	less, ok = xInfos.LessFileType(j, i)
	if ok {
		return less
	}

	return j < i
}

func sortInfoTypesDesc(items []os.FileInfo, compareDir fnCompareDir) {
	nameCachedInfos := infosTypesDesc{newInfosNames(items, compareDir)}
	sort.Sort(nameCachedInfos)
}

// sort size asc

type infosSizeAsc struct {
	infos
}

func (infos infosSizeAsc) Less(i, j int) bool {
	items := infos.items
	less, ok := infos.compareDir(items[i], items[j])
	if ok {
		return less
	}

	if items[i].Size() != items[j].Size() {
		return items[i].Size() < items[j].Size()
	}

	less, ok = util.CompareNumInFilename([]byte(items[i].Name()), []byte(items[j].Name()))
	if ok {
		return less
	}

	return i < j
}

func sortInfoSizesAsc(items []os.FileInfo, compareDir fnCompareDir) {
	infos := infosSizeAsc{newInfos(items, compareDir)}
	sort.Sort(infos)
}

// sort size desc

type infosSizeDesc struct {
	infos
}

func (infos infosSizeDesc) Less(i, j int) bool {
	items := infos.items
	less, ok := infos.compareDir(items[i], items[j])
	if ok {
		return less
	}

	if items[j].Size() != items[i].Size() {
		return items[j].Size() < items[i].Size()
	}

	less, ok = util.CompareNumInFilename([]byte(items[j].Name()), []byte(items[i].Name()))
	if ok {
		return less
	}

	return j < i
}

func sortInfoSizesDesc(items []os.FileInfo, compareDir fnCompareDir) {
	infos := infosSizeDesc{newInfos(items, compareDir)}
	sort.Sort(infos)
}

// sort time asc

type infosTimeAsc struct {
	infos
}

func (infos infosTimeAsc) Less(i, j int) bool {
	items := infos.items
	less, ok := infos.compareDir(items[i], items[j])
	if ok {
		return less
	}

	if !items[i].ModTime().Equal(items[j].ModTime()) {
		return items[i].ModTime().Before(items[j].ModTime())
	}

	less, ok = util.CompareNumInFilename([]byte(items[i].Name()), []byte(items[j].Name()))
	if ok {
		return less
	}

	return i < j
}

func sortInfoTimesAsc(items []os.FileInfo, compareDir fnCompareDir) {
	infos := infosTimeAsc{newInfos(items, compareDir)}
	sort.Sort(infos)
}

// sort time desc

type infosTimeDesc struct {
	infos
}

func (infos infosTimeDesc) Less(i, j int) bool {
	items := infos.items
	less, ok := infos.compareDir(items[i], items[j])
	if ok {
		return less
	}

	if !items[j].ModTime().Equal(items[i].ModTime()) {
		return items[j].ModTime().Before(items[i].ModTime())
	}

	less, ok = util.CompareNumInFilename([]byte(items[j].Name()), []byte(items[i].Name()))
	if ok {
		return less
	}

	return j < i
}

func sortInfoTimesDesc(items []os.FileInfo, compareDir fnCompareDir) {
	infos := infosTimeDesc{newInfos(items, compareDir)}
	sort.Sort(infos)
}

// sort original

type infosOriginalOrder struct {
	infos
}

func (infos infosOriginalOrder) Less(i, j int) bool {
	items := infos.items
	less, ok := infos.compareDir(items[i], items[j])
	if ok {
		return less
	}

	return i < j
}

func sortInfoOriginal(items []os.FileInfo, compareDir fnCompareDir) {
	infos := infosOriginalOrder{newInfos(items, compareDir)}
	sort.Sort(infos)
}

// sort

func sortInfos(items []os.FileInfo, sortBy string, defaultSortBy string) (sortInfo SortState) {
	if len(sortBy) == 0 {
		sortBy = defaultSortBy
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
		sortInfoNamesAsc(items, compareDir)
	case 'N':
		sortInfoNamesDesc(items, compareDir)
	case 'e':
		sortInfoTypesAsc(items, compareDir)
	case 'E':
		sortInfoTypesDesc(items, compareDir)
	case 's':
		sortInfoSizesAsc(items, compareDir)
	case 'S':
		sortInfoSizesDesc(items, compareDir)
	case 't':
		sortInfoTimesAsc(items, compareDir)
	case 'T':
		sortInfoTimesDesc(items, compareDir)
	default:
		if dirSort != dirSortMixed {
			sortInfoOriginal(items, compareDir)
		}
	}

	return SortState{dirSort, sortKey}
}
