package goNixArgParser

import (
	"os"
	"testing"
)

const testEnvVar = "GHFS_TEST_PARSE_RESULT"

func newTestResult(t *testing.T, hasEnv bool, env string, args []string) *ParseResult {
	t.Helper()

	// env is read while adding options, so set it before building the OptionSet
	if hasEnv {
		t.Setenv(testEnvVar, env)
	} else {
		os.Unsetenv(testEnvVar)
	}

	s := NewSimpleOptionSet()

	adds := []error{
		s.AddFlagValue("str", "--str", testEnvVar, "dft", ""),
		s.AddFlagValue("nodft", "--no-dft", "", "", ""),
		s.AddFlagValue("bool", "--bool", "", "true", ""),
		s.AddFlagValue("int", "--int", testEnvVar, "8080", ""),
		s.AddFlagValue("intnodft", "--int-no-dft", "", "", ""),
		s.AddFlagValues("multi", "--multi", "", []string{"d1"}, ""),
	}
	for _, err := range adds {
		if err != nil {
			t.Fatal(err)
		}
	}

	return s.Parse(args, nil)
}

// GetStringHasValue tells "key exists" from "value available": a source without a value is skipped
func TestGetStringHasValue(t *testing.T) {
	cases := []struct {
		name           string
		hasEnv         bool
		env            string
		args           []string
		key            string
		wantValue      string
		wantFoundKey   bool
		wantFoundValue bool
	}{
		{"flag with value", false, "", []string{"--str", "cli"}, "str", "cli", true, true},
		{"not specified falls to default", false, "", nil, "str", "dft", true, true},
		{"valueless flag falls to default", false, "", []string{"--str"}, "str", "dft", true, true},
		{"valueless flag falls to env", true, "env", []string{"--str"}, "str", "env", true, true},
		{"empty env falls to default", true, "", nil, "str", "dft", true, true},
		{"flag with value beats env", true, "env", []string{"--str", "cli"}, "str", "cli", true, true},
		{"valueless flag without default", false, "", []string{"--no-dft"}, "nodft", "", false, false},
		{"absent without default", false, "", nil, "nodft", "", false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newTestResult(t, c.hasEnv, c.env, c.args)

			value, foundKey, foundValue := r.GetStringHasValue(c.key)

			if value != c.wantValue || foundKey != c.wantFoundKey || foundValue != c.wantFoundValue {
				t.Errorf("got (%q, %v, %v), want (%q, %v, %v)",
					value, foundKey, foundValue, c.wantValue, c.wantFoundKey, c.wantFoundValue)
			}
		})
	}
}

// GetString keeps the original semantics: a present key wins even without a value,
// so a bare flag still clears the value from env or defaults
func TestGetStringKeepsValuelessKey(t *testing.T) {
	cases := []struct {
		name      string
		hasEnv    bool
		env       string
		args      []string
		wantValue string
		wantFound bool
	}{
		{"flag with value", false, "", []string{"--str", "cli"}, "cli", true},
		{"not specified falls to default", false, "", nil, "dft", true},
		{"valueless flag clears value", false, "", []string{"--str"}, "", true},
		{"valueless flag beats env", true, "env", []string{"--str"}, "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newTestResult(t, c.hasEnv, c.env, c.args)

			value, found := r.GetString("str")

			if value != c.wantValue || found != c.wantFound {
				t.Errorf("got (%q, %v), want (%q, %v)", value, found, c.wantValue, c.wantFound)
			}
		})
	}
}

// typed getters: an empty string never converts, so a source without a value is skipped
func TestGetBoolFallsBackToDefault(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantValue bool
		wantFound bool
	}{
		{"not specified uses default", nil, true, true},
		{"valueless flag uses default", []string{"--bool"}, true, true},
		{"explicit false wins", []string{"--bool", "false"}, false, true},
		{"explicit true wins", []string{"--bool", "true"}, true, true},
		{"unparsable value is not usable", []string{"--bool", "abc"}, false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newTestResult(t, false, "", c.args)

			value, found := r.GetBool("bool")

			if value != c.wantValue || found != c.wantFound {
				t.Errorf("got (%v, %v), want (%v, %v)", value, found, c.wantValue, c.wantFound)
			}
		})
	}
}

func TestGetIntFallsBackToDefault(t *testing.T) {
	cases := []struct {
		name      string
		hasEnv    bool
		env       string
		args      []string
		key       string
		wantValue int
		wantFound bool
	}{
		{"not specified uses default", false, "", nil, "int", 8080, true},
		{"valueless flag uses default", false, "", []string{"--int"}, "int", 8080, true},
		{"valueless flag uses env", true, "300", []string{"--int"}, "int", 300, true},
		{"explicit value wins", false, "", []string{"--int", "9090"}, "int", 9090, true},
		{"valueless flag without default", false, "", []string{"--int-no-dft"}, "intnodft", 0, false},
		{"unparsable value is not usable", false, "", []string{"--int", "abc"}, "int", 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newTestResult(t, c.hasEnv, c.env, c.args)

			value, found := r.GetInt(c.key)

			if value != c.wantValue || found != c.wantFound {
				t.Errorf("got (%v, %v), want (%v, %v)", value, found, c.wantValue, c.wantFound)
			}
		})
	}
}

// multi values are the opposite: a present key with no value is meaningful ("specified, no argument"),
// so neither skip the source nor degrade it to nil - downstream tells
// "disabled" from "enabled without argument" by nil versus empty slice
func TestGetStringsKeepsValuelessKey(t *testing.T) {
	t.Run("valueless flag yields non-nil empty slice", func(t *testing.T) {
		values, found := newTestResult(t, false, "", []string{"--multi"}).GetStrings("multi")

		if !found {
			t.Fatal("found should be true")
		}
		if values == nil {
			t.Fatal("values should not be nil")
		}
		if len(values) != 0 {
			t.Errorf("values should be empty, got %v", values)
		}
	})

	t.Run("not specified uses default", func(t *testing.T) {
		values, found := newTestResult(t, false, "", nil).GetStrings("multi")

		if !found {
			t.Fatal("found should be true")
		}
		if !expectStrings(values, "d1") {
			t.Errorf("got %v, want [d1]", values)
		}
	})

	t.Run("flag with values wins", func(t *testing.T) {
		values, found := newTestResult(t, false, "", []string{"--multi", "v1", "v2"}).GetStrings("multi")

		if !found {
			t.Fatal("found should be true")
		}
		if !expectStrings(values, "v1", "v2") {
			t.Errorf("got %v, want [v1 v2]", values)
		}
	})
}

// the last value wins when an option is provided multiple times
func TestGetValueUsesLastValue(t *testing.T) {
	newResult := func(t *testing.T, args []string) *ParseResult {
		t.Helper()

		s := NewSimpleOptionSet()

		adds := []error{
			s.AddFlagValue("single", "--single", "", "", ""),
			s.AddFlagValues("many", "--many", "", nil, ""),
			s.AddFlagValues("manyint", "--many-int", "", nil, ""),
		}
		for _, err := range adds {
			if err != nil {
				t.Fatal(err)
			}
		}

		return s.Parse(args, nil)
	}

	t.Run("multi values read as single value", func(t *testing.T) {
		value, found := newResult(t, []string{"--many", "v1", "v2", "v3"}).GetString("many")

		if !found {
			t.Fatal("found should be true")
		}
		if value != "v3" {
			t.Errorf("got %q, want %q", value, "v3")
		}
	})

	t.Run("option provided multiple times", func(t *testing.T) {
		value, found := newResult(t, []string{"--many", "v1", "--many", "v2"}).GetString("many")

		if !found {
			t.Fatal("found should be true")
		}
		if value != "v2" {
			t.Errorf("got %q, want %q", value, "v2")
		}
	})

	t.Run("typed getter uses last value too", func(t *testing.T) {
		value, found := newResult(t, []string{"--many-int", "1", "2", "3"}).GetInt("manyint")

		if !found {
			t.Fatal("found should be true")
		}
		if value != 3 {
			t.Errorf("got %v, want %v", value, 3)
		}
	})

	t.Run("single value option overrides previous", func(t *testing.T) {
		value, found := newResult(t, []string{"--single", "a", "--single", "b"}).GetString("single")

		if !found {
			t.Fatal("found should be true")
		}
		if value != "b" {
			t.Errorf("got %q, want %q", value, "b")
		}
	})
}
