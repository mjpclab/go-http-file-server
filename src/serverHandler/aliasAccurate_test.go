package serverHandler

import "testing"

func TestAliasAccurate(t *testing.T) {
	alias := createAliasAccurate("/hello/world/foo", "/tmp")

	// isMatch
	if !alias.isMatch("/hello/world/foo") {
		t.Error()
	}

	if alias.isMatch("/Hello/world/foo") {
		t.Error()
	}

	// isSuccessorOf
	if !alias.isSuccessorOf("/") {
		t.Error()
	}

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

	if alias.isSuccessorOf("/hi") {
		t.Error()
	}

	if alias.isSuccessorOf("/hi/") {
		t.Error()
	}

	if alias.isSuccessorOf("/hi/there") {
		t.Error()
	}

	// isPredecessorOf
	if !alias.isPredecessorOf("/hello/world/foo/bar") {
		t.Error()
	}

	if !alias.isPredecessorOf("/hello/world/foo/bar/") {
		t.Error()
	}

	if alias.isPredecessorOf("/hi/world/foo/bar") {
		t.Error()
	}
}
