package goNixArgParser

import "testing"

func TestRemoveEmpetyInplace(t *testing.T) {
	var output []string

	output = removeEmptyInplace([]string{"aa", "", "bb", "cc"})
	if !expectStrings(output, "aa", "bb", "cc") {
		t.Error(output)
	}
}

func TestRemoveQuotes(t *testing.T) {
	var output string

	output = removeQuotes("")
	if output != "" {
		t.Error(output)
	}

	output = removeQuotes("abc")
	if output != "abc" {
		t.Error(output)
	}

	output = removeQuotes(`hello "world" here`)
	if output != "hello world here" {
		t.Error(output)
	}

	output = removeQuotes(`hello 'world' there`)
	if output != "hello world there" {
		t.Error(output)
	}

	output = removeQuotes(`lorem 'ipsum`)
	if output != `lorem 'ipsum` {
		t.Error(output)
	}

	output = removeQuotes(`lorem "ipsum`)
	if output != `lorem "ipsum` {
		t.Error(output)
	}

	output = removeQuotes(`foo "bar" 'baz`)
	if output != `foo bar 'baz` {
		t.Error(output)
	}

	output = removeQuotes(`"Content-Security-Policy: default-src 'self'"`)
	if output != `Content-Security-Policy: default-src 'self'` {
		t.Error(output)
	}
}

func TestSplitLexicals(t *testing.T) {
	var output []string

	output = splitLexicals("")
	if !expectStrings(output) {
		t.Error(len(output), output)
	}

	output = splitLexicals(`  aaa bbb `)
	if !expectStrings(output, "aaa", "bbb") {
		t.Error(len(output), output)
	}

	output = splitLexicals(`aaa bbb "c c c" 'ddd' "ee'ee" 'ff"ff'		ggg`)
	if !expectStrings(output, "aaa", "bbb", `"c c c"`, "'ddd'", `"ee'ee"`, `'ff"ff'`, "ggg") {
		t.Error(len(output), output)
	}

	output = splitLexicals(`aa"bb"cc dd'ee'ff`)
	if !expectStrings(output, `aa"bb"cc`, `dd'ee'ff`) {
		t.Error(len(output), output)
	}

	// unpaired quotes
	output = splitLexicals(`aa"bb"cc"xxx dd'ee'ff'xxx`)
	if !expectStrings(output, `aa"bb"cc"xxx`, `dd'ee'ff'xxx`) {
		t.Error(len(output), output)
	}
}

func TestSplitToArgs(t *testing.T) {
	var output []string

	output = SplitToArgs(`  aaa bbb `)
	if !expectStrings(output, "aaa", "bbb") {
		t.Error(len(output), output)
	}

	output = SplitToArgs(`aaa bbb "c c c" 'ddd' "ee'ee" 'ff"ff'		ggg`)
	if !expectStrings(output, "aaa", "bbb", "c c c", "ddd", "ee'ee", `ff"ff`, "ggg") {
		t.Error(len(output), output)
	}

	output = SplitToArgs(`aa"bb"cc dd'ee'ff`)
	if !expectStrings(output, "aabbcc", "ddeeff") {
		t.Error(len(output), output)
	}

	// unpaired quotes
	output = SplitToArgs(`aa"bb"cc"xxx dd'ee'ff'xxx`)
	if !expectStrings(output, `aabbcc"xxx`, `ddeeff'xxx`) {
		t.Error(len(output), output)
	}
}
