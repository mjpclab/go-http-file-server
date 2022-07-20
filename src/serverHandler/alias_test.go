package serverHandler

import "testing"

func TestAlias(t *testing.T) {
	alias := createAlias("/hello/world/foo", "/tmp")

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

func TestAliasNextPartOf(t *testing.T) {
	var subName string
	var noMore, ok bool

	aliasAccurate := createAlias("/hello/world/foo", "/tmp")

	subName, noMore, ok = aliasAccurate.nextPartOf("/")
	if !ok {
		t.Error()
	}
	if subName != "hello" {
		t.Error()
	}
	if noMore {
		t.Error()
	}

	_, _, ok = aliasAccurate.nextPartOf("/test")
	if ok {
		t.Error()
	}

	subName, noMore, ok = aliasAccurate.nextPartOf("/hello")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if noMore {
		t.Error()
	}

	subName, noMore, ok = aliasAccurate.nextPartOf("/hello/")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if noMore {
		t.Error()
	}

	subName, noMore, ok = aliasAccurate.nextPartOf("/hello/world")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !noMore {
		t.Error()
	}

	subName, noMore, ok = aliasAccurate.nextPartOf("/hello/world/")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !noMore {
		t.Error()
	}

	subName, noMore, ok = aliasAccurate.nextPartOf("/hello/world/foo")
	if ok {
		t.Error()
	}

	subName, noMore, ok = aliasAccurate.nextPartOf("/hello/world/foo/")
	if ok {
		t.Error()
	}

	subName, noMore, ok = aliasAccurate.nextPartOf("/hello/world/foo/bar")
	if ok {
		t.Error()
	}

	subName, noMore, ok = aliasAccurate.nextPartOf("/hello/world/foo/bar/")
	if ok {
		t.Error()
	}
}
