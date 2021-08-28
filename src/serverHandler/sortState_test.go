package serverHandler

import (
	"testing"
)

func TestSortState_DirState(t *testing.T) {
	var sortBy string
	state := SortState{dirSortFirst, 'n'}

	sortBy = state.CurrentSort()
	if sortBy != "/n" {
		t.Error(sortBy)
	}

	sortBy = state.NextDirSort()
	if sortBy != "n/" {
		t.Error(sortBy)
	}

	state.dirSort = dirSortLast
	sortBy = state.NextDirSort()
	if sortBy != "n" {
		t.Error(sortBy)
	}

	state.dirSort = dirSortMixed
	sortBy = state.NextDirSort()
	if sortBy != "/n" {
		t.Error(sortBy)
	}
}

func TestSortState_KeyState(t *testing.T) {
	var sortBy string
	state := SortState{dirSortFirst, 'n'}

	sortBy = state.NextNameSort()
	if sortBy != "/N" {
		t.Error(sortBy)
	}

	state.key = 'N'
	sortBy = state.NextNameSort()
	if sortBy != "/n" {
		t.Error(sortBy)
	}

	sortBy = state.NextTypeSort()
	if sortBy != "/e" {
		t.Error(sortBy)
	}

	sortBy = state.NextSizeSort()
	if sortBy != "/S" {
		t.Error(sortBy)
	}

	sortBy = state.NextTimeSort()
	if sortBy != "/T" {
		t.Error(sortBy)
	}

	state.dirSort = dirSortMixed
	sortBy = state.NextTimeSort()
	if sortBy != "T" {
		t.Error(sortBy)
	}

	state.key = 'T'
	sortBy = state.NextTimeSort()
	if sortBy != "t" {
		t.Error(sortBy)
	}
}
