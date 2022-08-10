package param

import "testing"

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

func TestEntriesToHeaders(t *testing.T) {
	entries := []string{
		"",
		"key1:",
		":value2",
		"key3:value3",
	}
	headers := entriesToHeaders(entries)
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
