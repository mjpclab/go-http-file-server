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
	var maps map[string]string

	maps = normalizePathMaps([]string{":/data/lib://usr/lib"})
	if maps["/data/lib"] != "/usr/lib" {
		t.Error(maps)
	}

	maps = normalizePathMaps([]string{":/data/lib://usr/lib", "@foo@bar/baz"})
	if len(maps) != 2 {
		t.Error(maps)
	}
	if maps["/data/lib"] != "/usr/lib" {
		t.Error(maps)
	}
	if maps["/foo"] != "bar/baz" {
		t.Error(maps)
	}
}

func TestNormalizePathMapsNoCase(t *testing.T) {
	var maps map[string]string
	maps = normalizePathMapsNoCase([]string{":/data/lib://usr/lib"})
	if maps["/data/lib"] != "/usr/lib" {
		t.Error(maps)
	}

	maps = normalizePathMapsNoCase([]string{":/data/lib://usr/lib", "#/Data/Lib#/tmp/"})
	if len(maps) != 1 {
		t.Error(maps)
	}
	if maps["/Data/Lib"] != "/tmp" {
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
