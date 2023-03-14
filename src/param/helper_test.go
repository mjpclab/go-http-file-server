package param

import (
	"mjpclab.dev/ghfs/src/util"
	"path/filepath"
	"testing"
)

func TestEntriesToUsers(t *testing.T) {
	entries := []string{
		":pass1",
		"user2:",
		"user3:pass3",
	}
	users := entriesToUsers(entries)
	if len(users) != 3 {
		t.Fatal("user count is not 3")
	}
	if users[0][0] != "" {
		t.Fail()
	}
	if users[0][1] != "pass1" {
		t.Fail()
	}
	if users[1][0] != "user2" {
		t.Fail()
	}
	if users[1][1] != "" {
		t.Fail()
	}
	if users[2][0] != "user3" {
		t.Fail()
	}
	if users[2][1] != "pass3" {
		t.Fail()
	}
}

func TestNormalizePathRestrictAccesses(t *testing.T) {
	results, _ := normalizeAllPathValues([][]string{
		{"/foo", "host1", "host2"},
		{"/foo/", "host3", "host4"},
		{"/bar"},
	}, true, util.NormalizeUrlPath, util.ExtractHostsFromUrls)

	if len(results) != 2 {
		t.Error()
	}
	if !expectStrings(results[0], "/foo", "host1", "host2", "host3", "host4") {
		t.Error(results[0])
	}
	if !expectStrings(results[1], "/bar") {
		t.Error(results[1])
	}
}

func TestNormalizePathHeadersMap(t *testing.T) {
	var result [][]string

	result, _ = normalizeAllPathValues([][]string{
		{"/foo", "X-header1", "X-Value1"},
		{"/foo/", "X-header2", "X-Value2"},
		{"/bar", "X-header3", "X-Value3"},
		{"baz"},
		{"baz", ""},
		{"baz", "X-Not-Valid"},
		{"baz", "X-Not-Valid", ""},
	}, false, util.NormalizeUrlPath, normalizeHeaders)

	if len(result) != 2 {
		t.Error(result)
	}

	if !expectStrings(result[0], "/foo", "X-header1", "X-Value1", "X-header2", "X-Value2") {
		t.Error(result[0])
	}

	if !expectStrings(result[1], "/bar", "X-header3", "X-Value3") {
		t.Error(result[1])
	}
}

func TestNormalizePathMaps(t *testing.T) {
	var results [][2]string
	var fsPath string

	results, _ = normalizePathMaps([][2]string{{"/data/lib", "//usr/lib"}})
	if len(results) != 1 {
		t.Error(results)
	}
	fsPath, _ = filepath.Abs("/usr/lib")
	if !expectStrings(results[0][:], "/data/lib", fsPath) {
		t.Error(results[0])
	}

	results, _ = normalizePathMaps([][2]string{{"/data/lib", "//usr/lib"}, {"foo", "bar/baz"}})
	if len(results) != 2 {
		t.Error(results)
	}
	fsPath, _ = filepath.Abs("/usr/lib")
	if !expectStrings(results[0][:], "/data/lib", fsPath) {
		t.Error(results[0])
	}
	fsPath, _ = filepath.Abs("bar/baz")
	if !expectStrings(results[1][:], "/foo", fsPath) {
		t.Error(results[1])
	}

	results, _ = normalizePathMaps([][2]string{
		{"/data/lib", "//usr/lib"},
		{"foo", "bar/baz"},
		{"data/lib", "/usr/local/lib"},
	})
	if len(results) != 2 {
		t.Error(results)
	}
	fsPath, _ = filepath.Abs("/usr/local/lib")
	if !expectStrings(results[0][:], "/data/lib", fsPath) {
		t.Error(results[0])
	}
	fsPath, _ = filepath.Abs("bar/baz")
	if !expectStrings(results[1][:], "/foo", fsPath) {
		t.Error(results[1])
	}
}

func TestDedupPathValues(t *testing.T) {
	var result []string

	result = dedupPathValues(nil)
	if !expectStrings(result) {
		t.Error(result)
	}

	result = dedupPathValues([]string{})
	if !expectStrings(result) {
		t.Error(result)
	}

	result = dedupPathValues([]string{"/foo"})
	if !expectStrings(result, "/foo") {
		t.Error(result)
	}

	result = dedupPathValues([]string{"/foo", "wow"})
	if !expectStrings(result, "/foo", "wow") {
		t.Error(result)
	}

	result = dedupPathValues([]string{"/foo", "aa", "bb", "cc"})
	if !expectStrings(result, "/foo", "aa", "bb", "cc") {
		t.Error(result)
	}

	result = dedupPathValues([]string{"/foo", "xx", "yy", "xx", "zz", "zz", "/foo"})
	if !expectStrings(result, "/foo", "xx", "yy", "zz", "/foo") {
		t.Error(result)
	}
}

func TestNormalizeFilenames(t *testing.T) {
	files := []string{"", "abc/def.txt", "hello.txt"}
	normalized := normalizeFilenames(files)
	if len(normalized) != 1 || normalized[0] != "hello.txt" {
		t.Fail()
	}
}

func TestNormalizeHttpsPort(t *testing.T) {
	var httpsPort string
	var ok bool

	httpsPort, ok = normalizeToHttpsPort("123", []string{"123"})
	if !ok || httpsPort != ":123" {
		t.Error("1")
	}

	httpsPort, ok = normalizeToHttpsPort("234", []string{":234"})
	if !ok || httpsPort != ":234" {
		t.Error("2")
	}

	httpsPort, ok = normalizeToHttpsPort(":345", []string{"345"})
	if !ok || httpsPort != ":345" {
		t.Error("3")
	}

	httpsPort, ok = normalizeToHttpsPort(":456", []string{":456"})
	if !ok || httpsPort != ":456" {
		t.Error("4")
	}

	httpsPort, ok = normalizeToHttpsPort("", nil)
	if ok {
		t.Error("5")
	}

	httpsPort, ok = normalizeToHttpsPort("", []string{})
	if ok {
		t.Error("5")
	}

	httpsPort, ok = normalizeToHttpsPort("123", nil)
	if ok {
		t.Error("5")
	}

	httpsPort, ok = normalizeToHttpsPort("123", []string{})
	if ok {
		t.Error("5")
	}

	httpsPort, ok = normalizeToHttpsPort("", []string{""})
	if !ok || httpsPort != "" {
		t.Error("5")
	}

	httpsPort, ok = normalizeToHttpsPort("65536", []string{"65536"})
	if ok || httpsPort != "" {
		t.Error("6", httpsPort)
	}

	httpsPort, ok = normalizeToHttpsPort("", []string{"567"})
	if !ok || httpsPort != ":567" {
		t.Error("7")
	}

	httpsPort, ok = normalizeToHttpsPort("", []string{":678"})
	if !ok || httpsPort != ":678" {
		t.Error("8")
	}

	httpsPort, ok = normalizeToHttpsPort("789", []string{":890"})
	if ok || httpsPort != "" {
		t.Error("9")
	}

	httpsPort, ok = normalizeToHttpsPort("789", []string{"127.0.0.1:890"})
	if ok || httpsPort != "" {
		t.Error("10")
	}

	httpsPort, ok = normalizeToHttpsPort("789", []string{"[::1]:890"})
	if ok || httpsPort != "" {
		t.Error("11")
	}

	httpsPort, ok = normalizeToHttpsPort("", []string{":443"})
	if !ok || httpsPort != ":443" {
		t.Error("12")
	}

	httpsPort, ok = normalizeToHttpsPort("", []string{"127.0.0.1"})
	if !ok || httpsPort != "" {
		t.Error("13")
	}

	httpsPort, ok = normalizeToHttpsPort("", []string{"[::1]"})
	if !ok || httpsPort != "" {
		t.Error("14")
	}

	httpsPort, ok = normalizeToHttpsPort("443", []string{"127.0.0.1"})
	if !ok || httpsPort != "" {
		t.Error("15")
	}

	httpsPort, ok = normalizeToHttpsPort("443", []string{"[::1]"})
	if !ok || httpsPort != "" {
		t.Error("16")
	}
}
