package serverHandler

type dirSort int

const (
	dirSortFirst dirSort = -1
	dirSortMixed dirSort = 0
	dirSortLast  dirSort = 1
)

type SortState struct {
	dirSort dirSort
	key     byte
}

func (info SortState) DirSort() dirSort {
	return info.dirSort
}

func (info SortState) Key() string {
	return string(info.key)
}

func (info SortState) mergeDirWithKey(key byte) string {
	switch info.dirSort {
	case dirSortFirst:
		return "/" + string(key)
	case dirSortLast:
		return string(key) + "/"
	default:
		return string(key)
	}
}

func (info SortState) CurrentSort() string {
	return info.mergeDirWithKey(info.key)
}

func (info SortState) NextDirSort() string {
	switch info.dirSort {
	case dirSortFirst: // next is dirSortLast
		return string(info.key) + "/"
	case dirSortLast: // next is dirSortMixed
		return string(info.key)
	case dirSortMixed: // next is dirSortFirst
		return "/" + string(info.key)
	}
	return "/" + string(info.key)
}

func (info SortState) NextNameSort() string {
	var nextKey byte
	switch info.key {
	case 'n':
		nextKey = 'N'
	default:
		nextKey = 'n'
	}
	return info.mergeDirWithKey(nextKey)
}

func (info SortState) NextSizeSort() string {
	var nextKey byte
	switch info.key {
	case 's':
		nextKey = 'S'
	default:
		nextKey = 's'
	}
	return info.mergeDirWithKey(nextKey)
}

func (info SortState) NextTimeSort() string {
	var nextKey byte
	switch info.key {
	case 't':
		nextKey = 'T'
	default:
		nextKey = 't'
	}
	return info.mergeDirWithKey(nextKey)
}
