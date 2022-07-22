package param

import (
	"../util"
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

	_, _, k, v, ok = splitKeyValue("")
	if ok {
		t.Error("empty string should not OK")
	}

	_, _, k, v, ok = splitKeyValue(":")
	if ok {
		t.Error("separator-only string should not OK")
	}

	_, _, k, v, ok = splitKeyValue("::world")
	if ok {
		t.Error("empty key should not OK")
	}

	_, _, k, v, ok = splitKeyValue(":hello:")
	if ok {
		t.Error("empty value should not OK")
	}

	_, _, k, v, ok = splitKeyValue(":key:value")
	if !ok {
		t.Fail()
	}
	if k != "key" {
		t.Fail()
	}
	if v != "value" {
		t.Fail()
	}

	_, _, k, v, ok = splitKeyValue("@KEY@VALUE")
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

func TestNormalizePathRestrictAccesses(t *testing.T) {
	results, _ := normalizePathRestrictAccesses([]string{
		":/foo:host1:host2",
		":/foo/:host3:host4",
		":/bar",
	}, util.NormalizeUrlPath)

	if len(results) != 2 {
		t.Error()
	}
	if !expectStrings(results["/foo"], "host1", "host2", "host3", "host4") {
		t.Error()
	}
	if len(results["/bar"]) != 0 {
		t.Error()
	}
}

func TestNormalizePathHeadersMap(t *testing.T) {
	var result map[string][][2]string

	result, _ = normalizePathHeadersMap([]string{
		":/foo:X-header1:X-Value1",
		":/foo/:X-header2:X-Value2",
		":/bar:X-header3:X-Value3",
		":baz",
		":baz:",
		":baz:X-Not-Valid",
		":baz:X-Not-Valid:",
	}, util.NormalizeUrlPath)

	if len(result) != 2 {
		t.Error(result)
	}

	if len(result["/foo"]) != 2 {
		t.Error(result["/foo"])
	}
	if result["/foo"][0][0] != "X-header1" || result["/foo"][0][1] != "X-Value1" {
		t.Error(result["/foo"][0])
	}
	if result["/foo"][1][0] != "X-header2" || result["/foo"][1][1] != "X-Value2" {
		t.Error(result["/foo"][0])
	}

	if len(result["/bar"]) != 1 {
		t.Error(result["/foo"])
	}
	if result["/bar"][0][0] != "X-header3" || result["/bar"][0][1] != "X-Value3" {
		t.Error(result["/foo"][0])
	}
}

func TestNormalizePathMaps(t *testing.T) {
	var maps map[string]string
	var fsPath string

	maps, _ = normalizePathMaps([]string{":/data/lib://usr/lib"})
	fsPath, _ = filepath.Abs("/usr/lib")
	if maps["/data/lib"] != fsPath {
		t.Error(maps)
	}

	maps, _ = normalizePathMaps([]string{":/data/lib://usr/lib", "@foo@bar/baz"})
	if len(maps) != 2 {
		t.Error(maps)
	}
	if maps["/data/lib"] != "/usr/lib" {
		t.Error(maps)
	}
	fsPath, _ = filepath.Abs("bar/baz")
	if maps["/foo"] != fsPath {
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
