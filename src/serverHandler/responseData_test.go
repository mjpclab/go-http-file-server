package serverHandler

import "testing"

func TestGetPathEntries(t *testing.T) {
	var result []pathEntry

	result = getPathEntries("./", "/", true)
	if len(result) != 1 {
		t.Errorf("%#v\n", result)
	}
	if result[0].Path != "./" || result[0].Name != "/" {
		t.Error(result[0])
	}

	result = getPathEntries("./", "/a/b/c", false)
	if len(result) != 4 {
		t.Error(len(result))
	}
	if result[0].Path != "../../" || result[0].Name != "/" {
		t.Error(result[0])
	}
	if result[1].Path != "../" || result[1].Name != "a" {
		t.Error(result[1])
	}
	if result[2].Path != "./" || result[2].Name != "b" {
		t.Error(result[2])
	}
	if result[3].Path != "./c/" || result[3].Name != "c" {
		t.Error(result[3])
	}

	result = getPathEntries("./", "/a/b/c", true)
	if len(result) != 4 {
		t.Error(len(result))
	}
	if result[0].Path != "../../../" || result[0].Name != "/" {
		t.Error(result[0])
	}
	if result[1].Path != "../../" || result[1].Name != "a" {
		t.Error(result[1])
	}
	if result[2].Path != "../" || result[2].Name != "b" {
		t.Error(result[2])
	}
	if result[3].Path != "./" || result[3].Name != "c" {
		t.Error(result[3])
	}

	result = getPathEntries("./foo", "/", true)
	if len(result) != 1 {
		t.Error(len(result))
	}
	if result[0].Path != "./foo" || result[0].Name != "/" {
		t.Error(result[0])
	}
}
