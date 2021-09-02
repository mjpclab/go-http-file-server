package serverHandler

type dirSort int

const (
	dirSortFirst dirSort = -1
	dirSortMixed dirSort = 0
	dirSortLast  dirSort = 1
)

const (
	nameAsc  byte = 'n'
	nameDesc byte = 'N'
	typeAsc  byte = 'e'
	typeDesc byte = 'E'
	sizeAsc  byte = 's'
	sizeDesc byte = 'S'
	timeAsc  byte = 't'
	timeDesc byte = 'T'
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
	case nameAsc:
		nextKey = nameDesc
	default:
		nextKey = nameAsc
	}
	return info.mergeDirWithKey(nextKey)
}

func (info SortState) NextTypeSort() string {
	var nextKey byte
	switch info.key {
	case typeAsc:
		nextKey = typeDesc
	default:
		nextKey = typeAsc
	}
	return info.mergeDirWithKey(nextKey)
}

func (info SortState) NextSizeSort() string {
	var nextKey byte
	switch info.key {
	case sizeDesc:
		nextKey = sizeAsc
	default:
		nextKey = sizeDesc
	}
	return info.mergeDirWithKey(nextKey)
}

func (info SortState) NextTimeSort() string {
	var nextKey byte
	switch info.key {
	case timeDesc:
		nextKey = timeAsc
	default:
		nextKey = timeDesc
	}
	return info.mergeDirWithKey(nextKey)
}
