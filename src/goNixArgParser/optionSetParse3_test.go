package goNixArgParser

import (
	"testing"
)

func TestParse3(t *testing.T) {
	var err error

	s := NewOptionSet("", nil, []string{",,"})

	err = s.Append(&Option{
		Key:         "port",
		Flags:       []*Flag{{Name: "-p"}, {Name: "--port"}},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "root",
		Flags:       []*Flag{{Name: "-r"}, {Name: "--root"}},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	args := []string{
		"-p", "80",
		",,",
		"--port", "443", "--root", "/data/443",
	}

	configs := []string{
		"--root", "/data/80",
	}

	parsedGroups := s.ParseGroups(args, configs)
	if len(parsedGroups) != 2 {
		t.Fatal(len(parsedGroups))
	}

	parsed1 := parsedGroups[0]

	port1, _ := parsed1.GetString("port")
	if port1 != "80" {
		t.Error(port1)
	}

	root1, _ := parsed1.GetString("root")
	if root1 != "/data/80" {
		t.Error(root1)
	}

	parsed2 := parsedGroups[1]

	port2, _ := parsed2.GetString("port")
	if port2 != "443" {
		t.Error(port2)
	}

	root2, _ := parsed2.GetString("root")
	if root2 != "/data/443" {
		t.Error(root2)
	}
}
