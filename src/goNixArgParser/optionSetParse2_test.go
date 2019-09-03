package goNixArgParser

import (
	"fmt"
	"testing"
)

func TestParse2(t *testing.T) {
	var err error

	s := NewOptionSet("", nil)

	err = s.Append(&Option{
		Key:           "deft",
		Flags:         []*Flag{&Flag{Name: "-df"}, &Flag{Name: "--default"}},
		AcceptValue:   true,
		DefaultValues: []string{"myDefault"},
	})
	if err != nil {
		t.Error(err)
	}

	err = s.AddFlag("flag", "-flag", "flag option")
	if err != nil {
		t.Error(err)
	}

	err = s.AddFlag("p", "p", "flag p")
	if err != nil {
		t.Error(err)
	}

	err = s.AddFlag("q", "q", "flag q")
	if err != nil {
		t.Error(err)
	}

	err = s.AddFlags("flags", []string{"-flags", "--flags"}, "flags option")
	if err != nil {
		t.Error(err)
	}

	err = s.AddFlagValue("port", "--port", "", "21", "port to listen")
	if err != nil {
		t.Error(err)
	}

	err = s.AddFlagValues("ports", "--ports", "", []string{"80", "8080"}, "ports to listen for http")
	if err != nil {
		t.Error(err)
	}

	err = s.AddFlagsValue("file", []string{"-f", "--files"}, "", "", "file to open")
	if err != nil {
		t.Error(err)
	}

	err = s.AddFlagsValues("props", []string{"-p", "--props"}, "", []string{}, "properties")
	if err != nil {
		t.Error(err)
	}

	args := []string{
		"-flag",
		"pq",
		"--flags",
		"--port", "22",
		"--ports", "443", "4443",
		"--ports", "4444",
		"-f", "file1", "file2",
		"--props", "aa", "bb", "cc", "dd", "ee",
	}

	parsed := s.Parse(args)
	if v, _ := parsed.GetValue("deft"); v != "myDefault" {
		t.Error("deft")
	}
	if !parsed.HasFlagKey("flag") {
		t.Error("flag")
	}
	if !parsed.HasFlagKey("flags") {
		t.Error("flags")
	}
	if !parsed.HasFlagKey("p") {
		t.Error("p")
	}
	if !parsed.HasFlagKey("q") {
		t.Error("q")
	}
	if v, _ := parsed.GetValue("port"); v != "22" {
		t.Error("port")
	}

	ports, _ := parsed.GetValues("ports")
	fmt.Println("ports:", ports)
	if len(ports) != 3 {
		t.Error("ports")
	}

	if v, _ := parsed.GetValue("file"); v != "file1" {
		t.Error("file")
	}

	props, _ := parsed.GetValues("props")
	fmt.Println("props:", props)
	if len(props) != 5 {
		t.Error("props")
	}

	fmt.Println("rests:", parsed.GetRests())
	if len(parsed.GetRests()) != 1 {
		t.Error("rests")
	}
}
