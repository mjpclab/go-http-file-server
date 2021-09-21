package serverHandler

import "testing"

func TestGetAliasSubPartAccurate(t *testing.T) {
	var subName string
	var isLastPart, ok bool

	aliasAccurate := CreateAliasAccurate("/hello/world/foo", "/tmp")

	_, _, ok = getAliasSubPart(aliasAccurate, "/test")
	if ok {
		t.Error()
	}

	_, _, ok = getAliasSubPart(aliasAccurate, "/HELLO")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasAccurate, "/hello")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasAccurate, "/hello/")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasAccurate, "/hello/world")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasAccurate, "/hello/world/")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasAccurate, "/hello/world/foo")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasAccurate, "/hello/world/foo/")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasAccurate, "/hello/world/foo/bar")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasAccurate, "/hello/world/foo/bar/")
	if ok {
		t.Error()
	}

}

func TestGetAliasSubPartNoCase(t *testing.T) {
	var subName string
	var isLastPart, ok bool

	aliasNoCase := CreateAliasNoCase("/hello/world/foo", "/tmp")

	_, _, ok = getAliasSubPart(aliasNoCase, "/test")
	if ok {
		t.Error()
	}

	_, _, ok = getAliasSubPart(aliasNoCase, "/Test")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/HELLO")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello/")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/HELLO/")
	if !ok {
		t.Error()
	}
	if subName != "world" {
		t.Error()
	}
	if isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello/world")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/Hello/World/")
	if !ok {
		t.Error()
	}
	if subName != "foo" {
		t.Error()
	}
	if !isLastPart {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello/world/foo")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/Hello/world/foo")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello/world/foo/")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello/World/foo/")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello/world/foo/bar")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello/World/foo/bar")
	if ok {
		t.Error()
	}

	subName, isLastPart, ok = getAliasSubPart(aliasNoCase, "/hello/World/Foo/Bar/")
	if ok {
		t.Error()
	}

}
