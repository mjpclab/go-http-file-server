package serverHandler

import "testing"

func TestAliasNoCase(t *testing.T) {
	alias := createAliasNoCase("/hello/world/foo", "/tmp")

	// isMatch
	if !alias.isMatch("/hello/world/foo") {
		t.Error()
	}

	if !alias.isMatch("/Hello/world/foo") {
		t.Error()
	}

	// isSuccessorOf
	if !alias.isSuccessorOf("/") {
		t.Error()
	}

	if !alias.isSuccessorOf("/hello") {
		t.Error()
	}
	if !alias.isSuccessorOf("/Hello") {
		t.Error()
	}

	if !alias.isSuccessorOf("/hello/") {
		t.Error()
	}

	if !alias.isSuccessorOf("/HELLO/") {
		t.Error()
	}

	if !alias.isSuccessorOf("/hello/world") {
		t.Error()
	}

	if !alias.isSuccessorOf("/hello/world/") {
		t.Error()
	}

	if !alias.isSuccessorOf("/HELLO/world/") {
		t.Error()
	}

	if alias.isSuccessorOf("/hello/world/foo") {
		t.Error()
	}

	if alias.isSuccessorOf("/Hello/World/Foo/") {
		t.Error()
	}

	if alias.isSuccessorOf("/hello/world/foo/bar") {
		t.Error()
	}

	if alias.isSuccessorOf("/Hello/World/Foo/Bar/") {
		t.Error()
	}

	if alias.isSuccessorOf("/hi") {
		t.Error()
	}

	if alias.isSuccessorOf("/Hi/") {
		t.Error()
	}

	if alias.isSuccessorOf("/hi/There") {
		t.Error()
	}

	// isPredecessorOf
	if !alias.isPredecessorOf("/Hello/world/foo/bar") {
		t.Error()
	}

	if !alias.isPredecessorOf("/hello/world/FOO/BAR/") {
		t.Error()
	}

	if alias.isPredecessorOf("/Hi/world/foo/bar") {
		t.Error()
	}
}

func TestAliasSubPartNoCase(t *testing.T) {
	var subName string
	var isLastPart, ok bool

	aliasNoCase := createAliasNoCase("/hello/world/foo", "/tmp")

	subName, isLastPart, ok = aliasNoCase.subPart("/")
	if !ok {
		t.Error()
	}
	if subName != "hello" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	_, _, ok = aliasNoCase.subPart("/test")
	if ok {
		t.Error()
	}

	_, _, ok = aliasNoCase.subPart("/Test")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/HELLO")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello/")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/HELLO/")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello/world")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/Hello/World/")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello/world/foo")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/Hello/world/foo")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello/world/foo/")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello/World/foo/")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello/world/foo/bar")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello/World/foo/bar")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = aliasNoCase.subPart("/hello/World/Foo/Bar/")
	if ok {
		t.Error()
	}
}
