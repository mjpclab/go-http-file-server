package serverHandler

import (
	"os"
	"testing"
)

func TestGetMatchInfo(t *testing.T) {
	var match bool
	var childSel []string
	var info os.FileInfo

	var expect = func(isMatch bool, isChildSel ...string) bool {
		if isMatch != match {
			return false
		}

		if len(isChildSel) != len(childSel) {
			return false
		}

		if isChildSel != nil && childSel != nil {
			for i := 0; i < len(isChildSel); i++ {
				if isChildSel[i] != childSel[i] {
					return false
				}
			}
		}

		return true
	}

	info = createPlaceholderFileInfo("", true)
	match, childSel = matchSelection(info, nil)
	if !expect(true) {
		t.Error(match, childSel)
	}

	info = createPlaceholderFileInfo("", true)
	match, childSel = matchSelection(info, []string{})
	if !expect(true) {
		t.Error(match, childSel)
	}

	info = createPlaceholderFileInfo("", true)
	match, childSel = matchSelection(info, []string{"dir-x"})
	if !expect(false) {
		t.Error(match, childSel)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	match, childSel = matchSelection(info, nil)
	if !expect(true) {
		t.Error(match, childSel)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	match, childSel = matchSelection(info, []string{"dir-x"})
	if !expect(false) {
		t.Error(match, childSel)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	match, childSel = matchSelection(info, []string{"dir-a"})
	if !expect(true) {
		t.Error(match, childSel)
	}
	match, childSel = matchSelection(info, []string{"dir-a/"})
	if !expect(true) {
		t.Error(match, childSel)
	}

	info = createPlaceholderFileInfo("file-a", false)
	match, childSel = matchSelection(info, []string{"file-a"})
	if !expect(true) {
		t.Error(match, childSel)
	}
	match, childSel = matchSelection(info, []string{"file-a/"})
	if !expect(false) {
		t.Error(match, childSel)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	match, childSel = matchSelection(info, []string{"dir-a/dir-a1"})
	if !expect(true, "dir-a1") {
		t.Error(match, childSel)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	match, childSel = matchSelection(info, []string{"dir-a/dir-a1", "dir-a/dir-a2", "dir-a/dir-a1/dir-a11", "dir-b"})
	if !expect(true, "dir-a1", "dir-a2", "dir-a1/dir-a11") {
		t.Error(match, childSel)
	}
}
