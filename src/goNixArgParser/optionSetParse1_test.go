package goNixArgParser

import (
	"testing"
)

func TestParse1(t *testing.T) {
	var err error

	s := NewOptionSet("-", []string{"--"}, []string{",,"}, []string{"="}, []string{"-"})
	err = s.Add(Option{
		Key:         "tag",
		Summary:     "tag summary",
		Description: "tag description",
		Flags:       NewSimpleFlags([]string{"-t", "--tag"}),
		AcceptValue: false,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "single",
		Summary:     "single summary",
		Description: "single description",
		Flags:       NewSimpleFlags([]string{"-s", "--single"}),
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "multi",
		Flags:       NewSimpleFlags([]string{"-m", "--multi"}),
		AcceptValue: true,
		MultiValues: true,
		Delimiters:  []rune{','},
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:           "deft",
		Flags:         NewSimpleFlags([]string{"-df", "--default"}),
		AcceptValue:   true,
		DefaultValues: []string{"myDefault"},
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "singleMissingValue",
		Flags:       NewSimpleFlags([]string{"-sm", "--single-missing"}),
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "flagX",
		Flags:       NewSimpleFlags([]string{"-x"}),
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "flagY",
		Flags:       NewSimpleFlags([]string{"-y"}),
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:         "withEqual",
		Flags:       []*Flag{{Name: "--with-equal"}},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key: "withConcat",
		Flags: []*Flag{
			{Name: "-w", canConcatAssign: true},
		},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Add(Option{
		Key:           "fromenv",
		Flags:         []*Flag{{Name: "--from-env"}},
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
		"--single", "false",
		"xxx",
		"-m", "111", "222",
		"--multi", "333,444",
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
	r := s.Parse(args, nil)

	if r.HasFlagKey("deft") {
		t.Error("deft")
	}

	if v, _ := r.GetString("deft"); v != "myDefault" {
		t.Error("default")
	}

	r.SetConfigOption("deft", "cfgDefault")
	if deftValue, _ := r.GetString("deft"); deftValue != "cfgDefault" {
		t.Error(deftValue)
	}

	single, _ := r.GetString("single")
	if single != "false" {
		t.Error("single")
	}

	singleBool, _ := r.GetBool("single")
	if singleBool != false {
		t.Error(singleBool)
	}

	r.SetConfigOption("single", "cfg")
	if singleValue, _ := r.GetString("single"); singleValue != "false" {
		t.Error(singleValue)
	}

	multi, _ := r.GetStrings("multi")
	if len(multi) != 4 {
		t.Error("multi should have 4 values")
	}
	multiInts, _ := r.GetInts("multi")
	if len(multiInts) != 4 {
		t.Error("multiInts should have 4 values")
	}

	if !r.HasFlagKey("flagX") {
		t.Error("flagX")
	}

	if !r.HasFlagKey("flagY") {
		t.Error("flagY")
	}

	withEqual, _ := r.GetString("withEqual")
	if withEqual != "abcde" {
		t.Error("withEqual:", withEqual)
	}

	withConcat, _ := r.GetString("withConcat")
	if withConcat != "concatedvalue" {
		t.Error("withConcat:", withConcat)
	}

	// undefs: [-un1 --without-equal=bcdef -Wcannotconcat]
	undefs := r.GetUndefs()
	if len(undefs) != 3 {
		t.Error("undefs:", undefs)
	}
	if undefs[0] != "-un1" {
		t.Error(undefs[0])
	}
	if undefs[1] != "--without-equal" {
		t.Error(undefs[1])
	}
	if undefs[2] != "-Wcannotconcat" {
		t.Error(undefs[2])
	}
}
