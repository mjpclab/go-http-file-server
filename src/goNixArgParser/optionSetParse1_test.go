package goNixArgParser

import (
	"fmt"
	"testing"
)

func TestParse1(t *testing.T) {
	var err error

	s := NewOptionSet("-", []string{"--"})
	err = s.Append(&Option{
		Key:         "tag",
		Summary:     "tag summary",
		Description: "tag description",
		Flags:       NewSimpleFlags([]string{"-t", "--tag"}),
		AcceptValue: false,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "single",
		Summary:     "single summary",
		Description: "single description",
		Flags:       NewSimpleFlags([]string{"-s", "--single"}),
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "multi",
		Flags:       NewSimpleFlags([]string{"-m", "--multi"}),
		AcceptValue: true,
		MultiValues: true,
		Delimiters:  []rune{','},
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:           "deft",
		Flags:         NewSimpleFlags([]string{"-df", "--default"}),
		AcceptValue:   true,
		DefaultValues: []string{"myDefault"},
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "singleMissingValue",
		Flags:       NewSimpleFlags([]string{"-sm", "--single-missing"}),
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "flagX",
		Flags:       NewSimpleFlags([]string{"-x"}),
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "flagY",
		Flags:       NewSimpleFlags([]string{"-y"}),
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key: "withEqual",
		Flags: []*Flag{
			{Name: "--with-equal", canEqualAssign: true},
		},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key: "withConcat",
		Flags: []*Flag{
			{Name: "-w", canConcatAssign: true},
		},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key: "fromenv",
		Flags: []*Flag{
			{Name: "--from-env", canEqualAssign: true},
		},
		AcceptValue:   true,
		MultiValues:   true,
		Delimiters:    []rune{','},
		EnvVars:       []string{"FROMENV"},
		DefaultValues: []string{"fromEnvDft1", "fromEnvDft2"},
	})
	if err != nil {
		t.Error(err)
	}

	args := []string{
		"-t",
		"-un1", "val1",
		"--single", "singleval1",
		"xxx",
		"-m", "multival1", "multival2",
		"--multi", "multival3,multival4",
		"--with-equal=abcde",
		"--without-equal=bcdef",
		"-wconcatedvalue",
		"-Wcannotconcat",
		"-sq1", "'single_quoted_value'",
		"-sq2", "\"double_quoted_value\"",
		"-sm",
		"-xy",
		// "--from-env=1,2,3",
		"--",
		"-s", "1",
		"--with-equal=notwork",
	}
	r := s.Parse(args)
	fmt.Printf("%+v\n", r)

	if r.HasFlagKey("deft") {
		t.Error("deft")
	}

	if v, _ := r.GetValue("deft"); v != "myDefault" {
		t.Error("default")
	}

	single, _ := r.GetValue("single")
	fmt.Println("single:", single)
	if single != "singleval1" {
		t.Error("single")
	}

	multi, _ := r.GetValues("multi")
	fmt.Println("multi:", multi)
	if len(multi) != 4 {
		t.Error("multi should have 4 values")
	}

	if !r.HasFlagKey("flagX") {
		t.Error("flagX")
	}

	if !r.HasFlagKey("flagY") {
		t.Error("flagY")
	}

	withEqual, _ := r.GetValue("withEqual")
	if withEqual != "abcde" {
		t.Error("withEqual:", withEqual)
	}

	withConcat, _ := r.GetValue("withConcat")
	if withConcat != "concatedvalue" {
		t.Error("withConcat:", withConcat)
	}

	fmt.Println("rests:", r.rests)

	fmt.Print("fromenv: ")
	fmt.Println(r.GetValues("fromenv"))

	fmt.Print(string(s.GetHelp()))
}
