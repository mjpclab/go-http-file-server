package param

import (
	"mjpclab.dev/ghfs/src/util"
	"path/filepath"
	"testing"
)

func expectStrings(actuals []string, expects ...string) bool {
	if len(actuals) != len(expects) {
		return false
	}

	for i := range actuals {
		if actuals[i] != expects[i] {
			return false
		}
	}

	return true
}

func TestSplitKeyValues(t *testing.T) {
	var key string
	var values []string
	var ok bool

	key, values, ok = splitKeyValues("")
	if ok {
		t.Error()
	}

	key, values, ok = splitKeyValues(":")
	if ok {
		t.Error()
	}

	key, values, ok = splitKeyValues(":abc")
	if !ok {
		t.Error()
	}
	if key != "abc" {
		t.Error(key)
	}
	if len(values) != 0 {
		t.Error(values)
	}

	key, values, ok = splitKeyValues(":foo:")
	if !ok {
		t.Error()
	}
	if key != "foo" {
		t.Error(key)
	}
	if len(values) != 0 {
		t.Errorf("%#v\n", values)
	}

	key, values, ok = splitKeyValues(":foo:lorem:ipsum")
	if !ok {
		t.Error()
	}
	if key != "foo" {
		t.Error(key)
	}
	if !expectStrings(values, "lorem", "ipsum") {
		t.Error(values)
	}
}

func TestSplitKeyValue(t *testing.T) {
	var k, v string
	var ok bool

	k, v, ok = splitKeyValue("")
	if ok {
		t.Error("empty string should not OK")
	}

	k, v, ok = splitKeyValue(":")
	if ok {
		t.Error("separator-only string should not OK")
	}

	k, v, ok = splitKeyValue("::world")
	if ok {
		t.Error("empty key should not OK")
	}

	k, v, ok = splitKeyValue(":hello:")
	if ok {
		t.Error("empty value should not OK")
	}

	k, v, ok = splitKeyValue(":key:value")
	if !ok {
		t.Fail()
	}
	if k != "key" {
		t.Fail()
	}
	if v != "value" {
		t.Fail()
	}

	k, v, ok = splitKeyValue("@KEY@VALUE")
	if !ok {
		t.Fail()
	}
	if k != "KEY" {
		t.Fail()
	}
	if v != "VALUE" {
		t.Fail()
	}
}

func TestSplitAllKeyValue(t *testing.T) {
	results := splitAllKeyValue([]string{":foo:bar", "#lorem#ipsum"})
	if len(results) != 2 {
		t.Error(results)
	}
	if !expectStrings(results[0][:], "foo", "bar") {
		t.Error(results[0])
	}
	if !expectStrings(results[1][:], "lorem", "ipsum") {
		t.Error(results[1])
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

	httpsPort, ok = normalizeHttpsPort("123", []string{"123"})
	if !ok || httpsPort != ":123" {
		t.Error("1")
	}

	httpsPort, ok = normalizeHttpsPort("234", []string{":234"})
	if !ok || httpsPort != ":234" {
		t.Error("2")
	}

	httpsPort, ok = normalizeHttpsPort(":345", []string{"345"})
	if !ok || httpsPort != ":345" {
		t.Error("3")
	}

	httpsPort, ok = normalizeHttpsPort(":456", []string{":456"})
	if !ok || httpsPort != ":456" {
		t.Error("4")
	}

	httpsPort, ok = normalizeHttpsPort("", nil)
	if ok {
		t.Error("5")
	}

	httpsPort, ok = normalizeHttpsPort("", []string{})
	if ok {
		t.Error("5")
	}

	httpsPort, ok = normalizeHttpsPort("123", nil)
	if ok {
		t.Error("5")
	}

	httpsPort, ok = normalizeHttpsPort("123", []string{})
	if ok {
		t.Error("5")
	}

	httpsPort, ok = normalizeHttpsPort("", []string{""})
	if !ok || httpsPort != "" {
		t.Error("5")
	}

	httpsPort, ok = normalizeHttpsPort("65536", []string{"65536"})
	if ok || httpsPort != "" {
		t.Error("6", httpsPort)
	}

	httpsPort, ok = normalizeHttpsPort("", []string{"567"})
	if !ok || httpsPort != ":567" {
		t.Error("7")
	}

	httpsPort, ok = normalizeHttpsPort("", []string{":678"})
	if !ok || httpsPort != ":678" {
		t.Error("8")
	}

	httpsPort, ok = normalizeHttpsPort("789", []string{":890"})
	if ok || httpsPort != "" {
		t.Error("9")
	}

	httpsPort, ok = normalizeHttpsPort("789", []string{"127.0.0.1:890"})
	if ok || httpsPort != "" {
		t.Error("10")
	}

	httpsPort, ok = normalizeHttpsPort("789", []string{"[::1]:890"})
	if ok || httpsPort != "" {
		t.Error("11")
	}

	httpsPort, ok = normalizeHttpsPort("", []string{":443"})
	if !ok || httpsPort != ":443" {
		t.Error("12")
	}

	httpsPort, ok = normalizeHttpsPort("", []string{"127.0.0.1"})
	if !ok || httpsPort != "" {
		t.Error("13")
	}

	httpsPort, ok = normalizeHttpsPort("", []string{"[::1]"})
	if !ok || httpsPort != "" {
		t.Error("14")
	}

	httpsPort, ok = normalizeHttpsPort("443", []string{"127.0.0.1"})
	if !ok || httpsPort != "" {
		t.Error("15")
	}

	httpsPort, ok = normalizeHttpsPort("443", []string{"[::1]"})
	if !ok || httpsPort != "" {
		t.Error("16")
	}
}
