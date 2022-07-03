package goNixArgParser

import (
	"io"
	"strings"
)

func (opt *Option) isDelimiter(r rune) bool {
	for _, delimiter := range opt.Delimiters {
		if r == delimiter {
			return true
		}
	}
	return false
}

func (opt *Option) splitValues(input string) []string {
	values := strings.FieldsFunc(input, opt.isDelimiter)
	values = opt.filterValues(values)

	return values
}

func (opt *Option) filterValues(values []string) []string {
	if opt.UniqueValues {
		uniqueValues := make([]string, 0, len(values))
		uniqueValues = appendUnique(uniqueValues, values...)
		return uniqueValues
	}

	return values
}

func NewFlagOption(key, flag, envVar, summary string) Option {
	return Option{
		Key:     key,
		Flags:   []*Flag{NewSimpleFlag(flag)},
		EnvVars: stringToSlice(envVar),
		Summary: summary,
	}
}

func NewFlagsOption(key string, flags []string, envVar, summary string) Option {
	return Option{
		Key:     key,
		Flags:   NewSimpleFlags(flags),
		EnvVars: stringToSlice(envVar),
		Summary: summary,
	}
}

func NewFlagValueOption(key, flag, envVar, defaultValue, summary string) Option {
	return Option{
		Key:           key,
		Flags:         []*Flag{NewSimpleFlag(flag)},
		AcceptValue:   true,
		OverridePrev:  true,
		EnvVars:       stringToSlice(envVar),
		DefaultValues: stringToSlice(defaultValue),
		Summary:       summary,
	}
}

func NewFlagValuesOption(key, flag, envVar string, defaultValues []string, summary string) Option {
	return Option{
		Key:           key,
		Flags:         []*Flag{NewSimpleFlag(flag)},
		AcceptValue:   true,
		MultiValues:   true,
		UniqueValues:  true,
		EnvVars:       stringToSlice(envVar),
		DefaultValues: defaultValues,
		Summary:       summary,
	}
}

func NewFlagsValueOption(key string, flags []string, envVar, defaultValue, summary string) Option {
	return Option{
		Key:           key,
		Flags:         NewSimpleFlags(flags),
		AcceptValue:   true,
		OverridePrev:  true,
		EnvVars:       stringToSlice(envVar),
		DefaultValues: stringToSlice(defaultValue),
		Summary:       summary,
	}
}

func NewFlagsValuesOption(key string, flags []string, envVar string, defaultValues []string, summary string) Option {
	return Option{
		Key:           key,
		Flags:         NewSimpleFlags(flags),
		AcceptValue:   true,
		MultiValues:   true,
		UniqueValues:  true,
		EnvVars:       stringToSlice(envVar),
		DefaultValues: defaultValues,
		Summary:       summary,
	}
}

func (opt *Option) OutputHelp(w io.Writer) {
	if opt.Hidden {
		return
	}

	newline := []byte{'\n'}

	for i, flag := range opt.Flags {
		if i > 0 {
			w.Write([]byte{'|'})
		}
		io.WriteString(w, flag.Name)
	}

	if opt.AcceptValue {
		io.WriteString(w, " <value>")
		if opt.MultiValues {
			io.WriteString(w, " ...")
		}
	}

	w.Write(newline)

	if len(opt.EnvVars) > 0 {
		io.WriteString(w, "EnvVar: ")

		for i, envVar := range opt.EnvVars {
			if i > 0 {
				io.WriteString(w, ", ")
			}
			io.WriteString(w, envVar)
		}

		w.Write(newline)
	}

	if len(opt.DefaultValues) > 0 {
		io.WriteString(w, "Default: ")

		for i, d := range opt.DefaultValues {
			if i > 0 {
				io.WriteString(w, ", ")
			}
			io.WriteString(w, d)
		}

		w.Write(newline)
	}

	if len(opt.Summary) > 0 {
		io.WriteString(w, opt.Summary)
		w.Write(newline)
	}

	if len(opt.Description) > 0 {
		io.WriteString(w, opt.Description)
		w.Write(newline)
	}
}
