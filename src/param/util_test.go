package param

import "testing"

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

	key, values, ok = SplitKeyValues("")
	if ok {
		t.Error()
	}

	key, values, ok = SplitKeyValues(":")
	if ok {
		t.Error()
	}

	key, values, ok = SplitKeyValues(":abc")
	if !ok {
		t.Error()
	}
	if key != "abc" {
		t.Error(key)
	}
	if len(values) != 0 {
		t.Error(values)
	}

	key, values, ok = SplitKeyValues(":foo:")
	if !ok {
		t.Error()
	}
	if key != "foo" {
		t.Error(key)
	}
	if len(values) != 0 {
		t.Errorf("%#v\n", values)
	}

	key, values, ok = SplitKeyValues(":foo:lorem:ipsum")
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

	k, v, ok = SplitKeyValue("")
	if ok {
		t.Error("empty string should not OK")
	}

	k, v, ok = SplitKeyValue(":")
	if ok {
		t.Error("separator-only string should not OK")
	}

	k, v, ok = SplitKeyValue("::world")
	if ok {
		t.Error("empty key should not OK")
	}

	k, v, ok = SplitKeyValue(":hello:")
	if ok {
		t.Error("empty value should not OK")
	}

	k, v, ok = SplitKeyValue(":key:value")
	if !ok {
		t.Fail()
	}
	if k != "key" {
		t.Fail()
	}
	if v != "value" {
		t.Fail()
	}

	k, v, ok = SplitKeyValue("@KEY@VALUE")
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
	results := SplitAllKeyValue([]string{":foo:bar", "#lorem#ipsum"})
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

func TestEntriesToKVs(t *testing.T) {
	entries := []string{
		"",
		"key1:",
		":value2",
		"key3:value3",
	}
	headers := EntriesToKVs(entries)
	if len(headers) != 1 {
		t.Fatal("headers count should be 1", headers)
	}
	if headers[0][0] != "key3" {
		t.Error("key should be \"key3\"")
	}
	if headers[0][1] != "value3" {
		t.Error("value should be \"value3\"")
	}
}
