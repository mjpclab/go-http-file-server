package goNixArgParser

import (
	"testing"
)

func getGitCommand() *Command {
	cmdGit := NewSimpleCommand("git", "A version control tool")

	cmdGit.options.AddFlag("version", "--version", "", "display version")
	cmdGit.options.AddFlag("help", "--help", "", "show git help")

	cmdRemote := cmdGit.NewSimpleSubCommand("remote", "manage remotes", "rmt", "rt")
	cmdSetUrl := cmdRemote.NewSimpleSubCommand("set-url", "set remote url")
	cmdSetUrl.options.AddFlag("push", "--push", "", "")
	cmdSetUrl.options.AddFlagValue("dummy", "--dummy", "", "", "dummy option")
	cmdSetUrl.options.AddFlagValue("dummyX", "--dummy-x", "", "", "dummy-x option")

	cmdReset := cmdGit.NewSimpleSubCommand("reset", "reset command")
	cmdReset.options.AddFlag("hard", "--hard", "", "hard reset")
	cmdReset.options.AddFlag("mixed", "--mixed", "", "mixed reset")
	cmdReset.options.AddFlag("soft", "--soft", "", "soft reset")

	return cmdGit
}

/*
func TestGetLeafCmd(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "rmt", "set-url", "--push", "origin", "https://github.com/mjpclab/goNixArgParser.git"}
	_, normalizedArgs := cmd.getLeafCmd(args)
	for i, arg := range normalizedArgs {
		fmt.Printf("%d %+v\n", i, arg)
	}
}
*/

func TestParseCommand1(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "rmt", "set-url", "--push", "origin", "https://github.com/mjpclab/goNixArgParser.git"}

	result := cmd.Parse(args, nil)
	if result.commands[0] != "git" ||
		result.commands[1] != "remote" ||
		result.commands[2] != "set-url" {
		t.Error("commands", result.commands)
	}

	if !result.HasFlagKey("push") {
		t.Error("push")
	}

	if result.specifiedRests[0] != "origin" ||
		result.specifiedRests[1] != "https://github.com/mjpclab/goNixArgParser.git" {
		t.Error("rests", result.specifiedRests)
	}
}

func TestParseCommand2(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "xxx", "set-url", "origin", "https://github.com/mjpclab/goNixArgParser.git"}

	result := cmd.Parse(args, nil)
	if result.commands[0] != "git" ||
		result.commands[1] != "remote" {
		t.Error("commands", result.commands)
	}

	if result.specifiedRests[0] != "xxx" ||
		result.specifiedRests[1] != "set-url" ||
		result.specifiedRests[2] != "origin" ||
		result.specifiedRests[3] != "https://github.com/mjpclab/goNixArgParser.git" {
		t.Error("rests", result.specifiedRests)
	}
}

func TestParseCommand3(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "rmt", "set-url", "origin", "https://github.com/mjpclab/goNixArgParser.git"}
	configArgs := []string{"git", "rt", "set-url", "--dummy", "dummyconfigvalue"}
	result := cmd.Parse(args, configArgs)

	dummy, _ := result.GetString("dummy")
	if dummy != "dummyconfigvalue" {
		t.Error("dummy config value error", dummy)
	}

	configArgs = configArgs[1:]
	result = cmd.Parse(args, configArgs)

	dummy, _ = result.GetString("dummy")
	if dummy != "dummyconfigvalue" {
		t.Error("dummy config value error", dummy)
	}
}

func TestParseCommand4(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "--dummy", "dummy0", "github", "https://github.com/mjpclab/goNixArgParser.git", ",,", "--dummy", "dummy1", "bitbucket", "https://bitbucket.com/mjpclab/goNixArgParser.git"}
	configArgs := []string{"git", "remote", "set-url", "--dummy-x", "dummyXValue"}
	results := cmd.ParseGroups(args, configArgs)

	dummy0, _ := results[0].GetString("dummy")
	if dummy0 != "dummy0" {
		t.Error(results[0].GetStrings("dummy"))
	}

	dummyX, _ := results[0].GetString("dummyX")
	if dummyX != "dummyXValue" {
		t.Error(results[0].GetStrings("dummyX"))
	}

	dummy1, _ := results[1].GetString("dummy")
	if dummy1 != "dummy1" {
		t.Error(results[1].GetStrings("dummy"))
	}
}

func TestParseCommand5(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "--dummy", "dummy0"}
	configArgs := []string{"git", "no-such-cmd", "set-url", "--dummy-x", "dummyXValue"}
	result := cmd.Parse(args, configArgs)

	dummyX, _ := result.GetString("dummyX")
	if dummyX != "" {
		t.Error("dummyX")
	}
}

func TestParseCommand6(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "github", "https://github.com/mjpclab/goNixArgParser.git", ",,", "bitbucket", "https://bitbucket.com/mjpclab/goNixArgParser.git"}
	configArgs := []string{"git", "remote", "set-url", "--dummy", "dummy0", ",,", "--dummy", "dummy1"}
	results := cmd.ParseGroups(args, configArgs)

	dummy0, _ := results[0].GetString("dummy")
	if dummy0 != "dummy0" {
		t.Error(results[0].GetStrings("dummy"))
	}

	dummy1, _ := results[1].GetString("dummy")
	if dummy1 != "dummy1" {
		t.Error(results[1].GetStrings("dummy"))
	}

	configArgs = configArgs[1:]
	results = cmd.ParseGroups(args, configArgs)

	dummy0, _ = results[0].GetString("dummy")
	if dummy0 != "dummy0" {
		t.Error(results[0].GetStrings("dummy"))
	}

	dummy1, _ = results[1].GetString("dummy")
	if dummy1 != "dummy1" {
		t.Error(results[1].GetStrings("dummy"))
	}
}

func TestParseCommand7(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "github", "https://github.com/mjpclab/goNixArgParser.git"}
	configArgs := []string{"--dummy", "dummy0"}

	result := cmd.Parse(args, configArgs)

	cmdPaths := result.GetCommands()
	if len(cmdPaths) != 3 {
		t.Error(len(cmdPaths))
	}
	if cmdPaths[0] != "git" {
		t.Error(cmdPaths[0])
	}
	if cmdPaths[1] != "remote" {
		t.Error(cmdPaths[1])
	}
	if cmdPaths[2] != "set-url" {
		t.Error(cmdPaths[2])
	}

	dummy, _ := result.GetString("dummy")
	if dummy != "dummy0" {
		t.Error(dummy)
	}
}
