package goNixArgParser

import (
	"testing"
)

func TestParse3(t *testing.T) {
	var err error

	s := NewOptionSet("-", nil, []string{",,"}, []string{"="}, []string{"-"})

	err = s.Add(Option{
		Key:         "bool",
		Flags:       []*Flag{{Name: "-b", canMerge: true}, {Name: "--bool"}},
		AcceptValue: false,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "port",
		Flags:       []*Flag{{Name: "-p", canMerge: true, canFollowAssign: true}, {Name: "--port", canFollowAssign: true}},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "root",
		Flags:       []*Flag{{Name: "-r"}, {Name: "--root", canFollowAssign: true}},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "username",
		Flags:       []*Flag{{Name: "--username", prefixMatchLen: 1, canFollowAssign: true}},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "password",
		Flags:       []*Flag{{Name: "--password", prefixMatchLen: 1}},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	args := []string{
		"-p", "80",
		",,",
		"--port", "443", "--root", "/data/443",
		"--user", "user1",
		"--pass=pass1",
		",,",
	}

	configs := []string{
		"--root", "/data/80",
	}

	// groups
	parsedGroups := s.ParseGroups(args, configs)
	if len(parsedGroups) != 3 {
		t.Error("parsed group count:", len(parsedGroups))
	}

	// groups - group 1
	parsed1 := parsedGroups[0]

	port1, _ := parsed1.GetString("port")
	if port1 != "80" {
		t.Error(port1)
	}

	root1, _ := parsed1.GetString("root")
	if root1 != "/data/80" {
		t.Error(root1)
	}

	// groups - group 2
	parsed2 := parsedGroups[1]

	port2, _ := parsed2.GetString("port")
	if port2 != "443" {
		t.Error(port2)
	}

	root2, _ := parsed2.GetString("root")
	if root2 != "/data/443" {
		t.Error(root2)
	}

	username, _ := parsed2.GetString("username")
	if username != "user1" {
		t.Error(username)
	}

	password, _ := parsed2.GetString("password")
	if password != "pass1" {
		t.Error(password)
	}

	// merge & equal assign
	args = []string{"-bp=8080"}
	result := s.Parse(args, nil)
	if !result.HasKey("bool") {
		t.Error("bool")
	}

	port, _ := result.GetString("port")
	if port != "8080" {
		t.Error(port)
	}
}
