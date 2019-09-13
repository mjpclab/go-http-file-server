package goNixArgParser

import (
	"fmt"
	"testing"
)

func getGitCommand() *Command {
	cmdGit := NewSimpleCommand("git", "A version control tool")

	cmdGit.OptionSet.AddFlag("version", "--version", "", "display version")
	cmdGit.OptionSet.AddFlag("help", "--help", "", "show git help")

	cmdSetUrl := cmdGit.NewSimpleSubCommand("remote", "manage remotes").NewSimpleSubCommand("set-url", "set remote url")
	cmdSetUrl.OptionSet.AddFlag("push", "--push", "", "")

	cmdReset := cmdGit.NewSimpleSubCommand("reset", "reset command")
	cmdReset.OptionSet.AddFlag("hard", "--hard", "", "hard reset")
	cmdReset.OptionSet.AddFlag("mixed", "--mixed", "", "mixed reset")
	cmdReset.OptionSet.AddFlag("soft", "--soft", "", "soft reset")

	return cmdGit
}

func TestNormalizeCmdArgs(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "--push", "origin", "https://github.com/mjpclab/goNixArgParser.git"}
	normalizedArgs, _ := cmd.getNormalizedArgs(args)
	for i, arg := range normalizedArgs {
		fmt.Printf("%d %+v\n", i, arg)
	}
}

func TestParseCommand1(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "--push", "origin", "https://github.com/mjpclab/goNixArgParser.git"}

	result := cmd.Parse(args, nil)
	if result.commands[0] != "git" ||
		result.commands[1] != "remote" ||
		result.commands[2] != "set-url" {
		t.Error("commands", result.commands)
	}

	if !result.HasFlagKey("push") {
		t.Error("push")
	}

	if result.argRests[0] != "origin" ||
		result.argRests[1] != "https://github.com/mjpclab/goNixArgParser.git" {
		t.Error("rests", result.argRests)
	}

	cmd.PrintHelp()
}

func TestParseCommand2(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "xxx", "set-url", "origin", "https://github.com/mjpclab/goNixArgParser.git"}

	result := cmd.Parse(args, nil)
	if result.commands[0] != "git" ||
		result.commands[1] != "remote" {
		t.Error("commands", result.commands)
	}

	if result.argRests[0] != "xxx" ||
		result.argRests[1] != "set-url" ||
		result.argRests[2] != "origin" ||
		result.argRests[3] != "https://github.com/mjpclab/goNixArgParser.git" {
		t.Error("rests", result.argRests)
	}
}
