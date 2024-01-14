package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
	"testing"
)

func TestPathStrings(t *testing.T) {
	ps := pathStringsList{
		{"/a", []string{"a"}},
		{"/a/b", []string{"ab"}},
		{"/a/b/c", []string{"abc"}},
		{"/foo/bar", []string{"foobar"}},
	}

	mergeWith := []string{"/xxx", "/yyy"}
	merged := ps.mergePrefixMatched(mergeWith, util.HasUrlPrefixDir, "/a/b")
	if len(mergeWith) != 2 {
		t.Error()
	}
	if len(merged) != 4 || merged[2] != "a" || merged[3] != "ab" {
		t.Error(merged)
	}

	merged = ps.mergePrefixMatched(nil, util.HasUrlPrefixDir, "/lorem/ipsum")
	if merged != nil {
		t.Error(merged)
	}

	successors := ps.filterSuccessor(util.HasUrlPrefixDir, "/a/b")
	if len(successors) != 1 || successors[0].path != "/a/b/c" {
		t.Error(successors)
	}
}
