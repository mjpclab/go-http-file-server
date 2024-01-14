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

func TestPathHeaders(t *testing.T) {
	ps := pathHeadersList{
		{"/a", [][2]string{{"a", "a-value"}}},
		{"/a/b", [][2]string{{"ab", "ab-value"}}},
		{"/a/b/c", [][2]string{{"abc", "abc-value"}}},
		{"/foo/bar", [][2]string{{"foobar", "foobar-value"}}},
	}

	mergeWith := [][2]string{{"Access-Control-Allow-Origin", "*"}, {"Access-Control-Allow-Headers", "*"}}
	merged := ps.mergePrefixMatched(mergeWith, util.HasUrlPrefixDir, "/a/b")
	if len(mergeWith) != 2 {
		t.Error()
	}
	if len(merged) != 4 || merged[2][0] != "a" || merged[3][0] != "ab" {
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

func TestPathList(t *testing.T) {
	list := []string{
		"/a", "/a/b", "/a/b/c", "/foo/bar",
	}

	if !prefixMatched(list, util.HasUrlPrefixDir, "/a/b") {
		t.Error()
	}
	if !prefixMatched(list, util.HasUrlPrefixDir, "/a/b/c/d") {
		t.Error()
	}
	if prefixMatched(list, util.HasUrlPrefixDir, "/lorem/ipsum") {
		t.Error()
	}

	successors := filterSuccessor(list, util.HasUrlPrefixDir, "/a/b")
	if len(successors) != 1 || successors[0] != "/a/b/c" {
		t.Error(successors)
	}
}
