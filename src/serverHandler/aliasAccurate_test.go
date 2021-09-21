package serverHandler

import "testing"

func TestAliasAccurate(t *testing.T) {
	alias := CreateAliasAccurate("/hello/world/foo", "/tmp")

	// isMatch
	if !alias.isMatch("/hello/world/foo") {
		t.Error()
	}

	if alias.isMatch("/Hello/world/foo") {
		t.Error()
	}

	// isSuccessorOf
	if !alias.isSuccessorOf("/hello") {
		t.Error()
	}

	if !alias.isSuccessorOf("/hello/") {
		t.Error()
	}

	if !alias.isSuccessorOf("/hello/world") {
		t.Error()
	}

	if !alias.isSuccessorOf("/hello/world/") {
		t.Error()
	}

	if alias.isSuccessorOf("/HELLO/world/") {
		t.Error()
	}

	if alias.isSuccessorOf("/hello/world/foo") {
		t.Error()
	}

	if alias.isSuccessorOf("/hello/world/foo/") {
		t.Error()
	}

	if alias.isSuccessorOf("/hello/world/foo/bar") {
		t.Error()
	}

	if alias.isSuccessorOf("/hello/world/foo/bar/") {
		t.Error()
	}
}
