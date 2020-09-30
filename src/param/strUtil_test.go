package param

import "testing"

func TestSplitMapping(t *testing.T) {
	var k, v string
	var ok bool

	k, v, ok = splitMapping("")
	if ok {
		t.Error("empty string should not OK")
	}

	k, v, ok = splitMapping(":")
	if ok {
		t.Error("separator-only string should not OK")
	}

	k, v, ok = splitMapping("::world")
	if ok {
		t.Error("empty key should not OK")
	}

	k, v, ok = splitMapping(":hello:")
	if ok {
		t.Error("empty value should not OK")
	}

	k, v, ok = splitMapping(":key:value")
	if !ok {
		t.Fail()
	}
	if k != "key" {
		t.Fail()
	}
	if v != "value" {
		t.Fail()
	}
}

func TestNormalizePathMaps(t *testing.T) {
	maps := normalizePathMaps([]string{":/data/lib://usr/lib"})
	if maps["/data/lib"] != "/usr/lib" {
		t.Error(maps)
	}
}

func TestNormalizeFilenames(t *testing.T) {
	files := []string{"", "abc/def.txt", "hello.txt"}
	normalized := normalizeFilenames(files)
	if len(normalized) != 1 || normalized[0] != "hello.txt" {
		t.Fail()
	}
}
