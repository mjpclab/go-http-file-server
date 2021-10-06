package serverHandler

import (
	"os"
	"testing"
)

func TestGetMatchInfo(t *testing.T) {
	var matchName, matchPrefix bool
	var childList []string
	var info os.FileInfo

	var expect = func(isMatchName, isMatchPrefix bool, isChildList ...string) bool {
		if isMatchName != matchName {
			return false
		}
		if isMatchPrefix != matchPrefix {
			return false
		}

		if len(isChildList) != len(childList) {
			return false
		}

		if len(isChildList) != len(childList) {
			return false
		}

		if isChildList != nil && childList != nil {
			for i := 0; i < len(isChildList); i++ {
				if isChildList[i] != childList[i] {
					return false
				}
			}
		}

		return true
	}

	info = createPlaceholderFileInfo("", true)
	matchName, matchPrefix, childList = matchSelection(info, nil)
	if !expect(true, false) {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfo("", true)
	matchName, matchPrefix, childList = matchSelection(info, []string{})
	if !expect(true, false) {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfo("", true)
	matchName, matchPrefix, childList = matchSelection(info, []string{"dir-x"})
	if !expect(false, false) {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	matchName, matchPrefix, childList = matchSelection(info, nil)
	if !expect(true, false) {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	matchName, matchPrefix, childList = matchSelection(info, []string{"dir-x"})
	if !expect(false, false) {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	matchName, matchPrefix, childList = matchSelection(info, []string{"dir-a"})
	if !expect(true, false) {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	matchName, matchPrefix, childList = matchSelection(info, []string{"dir-a/dir-a1"})
	if !expect(false, true, "dir-a1") {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfo("dir-a", true)
	matchName, matchPrefix, childList = matchSelection(info, []string{"dir-a/dir-a1", "dir-a/dir-a2", "dir-a/dir-a1/dir-a11", "dir-b"})
	if !expect(false, true, "dir-a1", "dir-a2", "dir-a1/dir-a11") {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfoNoCase("dir-a", true)
	matchName, matchPrefix, childList = matchSelection(info, []string{"Dir-a/dir-a1"})
	if !expect(false, true, "dir-a1") {
		t.Error(matchName, matchPrefix, childList)
	}

	info = createPlaceholderFileInfoNoCase("dir-a", true)
	matchName, matchPrefix, childList = matchSelection(info, []string{"Dir-a/dir-a1", "dir-a/dir-a2", "dir-a/dir-a1/dir-a11", "dir-b"})
	if !expect(false, true, "dir-a1", "dir-a2", "dir-a1/dir-a11") {
		t.Error(matchName, matchPrefix, childList)
	}
}
