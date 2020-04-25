package serverHandler

import "html/template"

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

func (info SortState) mergeDirWithKey(key byte) template.HTML {
	switch info.dirSort {
	case dirSortFirst:
		return "/" + template.HTML(key)
	case dirSortLast:
		return template.HTML(key) + "/"
	default:
		return template.HTML(key)
	}
}

func (info SortState) CurrentSort() template.HTML {
	return info.mergeDirWithKey(info.key)
}

func (info SortState) NextDirSort() template.HTML {
	switch info.dirSort {
	case dirSortFirst: // next is dirSortLast
		return template.HTML(info.key) + "/"
	case dirSortLast: // next is dirSortMixed
		return template.HTML(info.key)
	case dirSortMixed: // next is dirSortFirst
		return "/" + template.HTML(info.key)
	}
	return "/" + template.HTML(info.key)
}

func (info SortState) NextNameSort() template.HTML {
	var nextKey byte
	switch info.key {
	case 'n':
		nextKey = 'N'
	default:
		nextKey = 'n'
	}
	return info.mergeDirWithKey(nextKey)
}

func (info SortState) NextSizeSort() template.HTML {
	var nextKey byte
	switch info.key {
	case 's':
		nextKey = 'S'
	default:
		nextKey = 's'
	}
	return info.mergeDirWithKey(nextKey)
}

func (info SortState) NextTimeSort() template.HTML {
	var nextKey byte
	switch info.key {
	case 't':
		nextKey = 'T'
	default:
		nextKey = 't'
	}
	return info.mergeDirWithKey(nextKey)
}

func (info SortState) Key() template.HTML {
	return template.HTML(info.key)
}
