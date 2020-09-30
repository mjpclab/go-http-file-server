package param

import "testing"

func TestEntriesToUsers(t *testing.T) {
	entries := []string{
		":pass1",
		"user2:",
		"user3:pass3",
	}
	users := EntriesToUsers(entries)
	if len(users) != 3 {
		t.Fatal("user count is not 3")
	}
	if users[0].Username != "" {
		t.Fail()
	}
	if users[0].Password != "pass1" {
		t.Fail()
	}
	if users[1].Username != "user2" {
		t.Fail()
	}
	if users[1].Password != "" {
		t.Fail()
	}
	if users[2].Username != "user3" {
		t.Fail()
	}
	if users[2].Password != "pass3" {
		t.Fail()
	}
}
