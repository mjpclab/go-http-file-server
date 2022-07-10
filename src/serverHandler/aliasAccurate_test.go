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

func TestAliasSubPartAccurate(t *testing.T) {
	var subName string
	var isLastPart, ok bool

	aliasAccurate := createAliasAccurate("/hello/world/foo", "/tmp")

	subName, isLastPart, ok = aliasAccurate.subPart("/")
	if !ok {
		t.Error()
	}
	if subName != "hello" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	_, _, ok = aliasAccurate.subPart("/test")
	if ok {
		t.Error()
	}

	_, _, ok = aliasAccurate.subPart("/HELLO")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasAccurate.subPart("/hello")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasAccurate.subPart("/hello/")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasAccurate.subPart("/hello/world")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasAccurate.subPart("/hello/world/")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasAccurate.subPart("/hello/world/foo")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasAccurate.subPart("/hello/world/foo/")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasAccurate.subPart("/hello/world/foo/bar")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasAccurate.subPart("/hello/world/foo/bar/")
	if ok {
		t.Error()
	}
}
