package goNixArgParser

import "os"

var CommandLine *OptionSet = NewOptionSet("-")

func Append(opt *Option) error {
	return CommandLine.Append(opt)
}

func AddFlag(key, flag, summary string) error {
	return CommandLine.AddFlag(key, flag, summary)
}

func AddFlags(key string, flags []string, summary string) error {
	return CommandLine.AddFlags(key, flags, summary)
}

func AddFlagValue(key, flag, defaultValue, summary string) error {
	return CommandLine.AddFlagValue(key, flag, defaultValue, summary)
}

func AddFlagValues(key, flag string, defaultValues []string, summary string) error {
	return CommandLine.AddFlagValues(key, flag, defaultValues, summary)
}

func AddFlagsValue(key string, flags []string, defaultValue, summary string) error {
	return CommandLine.AddFlagsValue(key, flags, defaultValue, summary)
}

func AddFlagsValues(key string, flags, defaultValues []string, summary string) error {
	return CommandLine.AddFlagsValues(key, flags, defaultValues, summary)
}

func PrintHelp() {
	CommandLine.PrintHelp()
}

func Parse() *ParseResult {
	return CommandLine.Parse(os.Args[1:])
}
