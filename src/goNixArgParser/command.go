package goNixArgParser

import (
	"bytes"
	"os"
	"path"
)

func NewCommand(
	name, summary, mergeFlagPrefix string,
	restsSigns, groupSeps []string,
) *Command {
	return &Command{
		name:        name,
		summary:     summary,
		options:     NewOptionSet(mergeFlagPrefix, restsSigns, groupSeps),
		subCommands: []*Command{},
	}
}

func NewSimpleCommand(name, summary string) *Command {
	return &Command{
		name:        name,
		summary:     summary,
		options:     NewSimpleOptionSet(),
		subCommands: []*Command{},
	}
}

func (c *Command) NewSubCommand(
	name, summary, mergeFlagPrefix string,
	restsSigns, groupSeps []string,
) *Command {
	subCommand := NewCommand(name, summary, mergeFlagPrefix, restsSigns, groupSeps)
	c.subCommands = append(c.subCommands, subCommand)
	return subCommand
}

func (c *Command) NewSimpleSubCommand(name, summary string) *Command {
	subCommand := NewSimpleCommand(name, summary)
	c.subCommands = append(c.subCommands, subCommand)
	return subCommand
}

func (c *Command) GetSubCommand(name string) *Command {
	if c.subCommands == nil {
		return nil
	}

	for _, cmd := range c.subCommands {
		if cmd.name == name {
			return cmd
		}
	}
	return nil
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Summary() string {
	return c.summary
}

func (c *Command) Options() *OptionSet {
	return c.options
}

func (c *Command) SubCommands() []*Command {
	return c.subCommands
}

func (c *Command) getNormalizedArgs(initArgs []string) (*Command, []*Arg) {
	cmd := c

	if len(initArgs) == 0 {
		return cmd, []*Arg{}
	}

	args := make([]*Arg, 0, len(initArgs))

	for i, arg := range initArgs {
		if i == 0 && cmd.name == arg {
			args = append(args, NewArg(arg, CommandArg))
		} else if subCmd := cmd.GetSubCommand(arg); subCmd != nil {
			args = append(args, NewArg(arg, CommandArg))
			cmd = subCmd
		} else {
			break
		}
	}

	return cmd, args
}

func (c *Command) splitCommandsArgs(initArgs, initConfigs []string) (
	argsLeafCmd *Command,
	commands, optionSetInitArgs, optionSetInitConfigs []string,
) {
	argsLeafCmd, argCmds := c.getNormalizedArgs(initArgs)
	configsLeafCmd, configCmds := c.getNormalizedArgs(initConfigs)

	commands = []string{}
	for _, arg := range argCmds {
		if arg.Type != CommandArg {
			break
		}
		commands = append(commands, arg.Text)
	}

	optionSetInitArgs = initArgs[len(argCmds):]

	if argsLeafCmd == configsLeafCmd {
		optionSetInitConfigs = initConfigs[len(configCmds):]
	} else {
		optionSetInitConfigs = []string{}
	}

	return
}

func (c *Command) Parse(initArgs, initConfigs []string) *ParseResult {
	leafCmd, commands, optionSetInitArgs, optionSetInitConfigs := c.splitCommandsArgs(initArgs, initConfigs)
	result := leafCmd.options.Parse(optionSetInitArgs, optionSetInitConfigs)
	result.commands = commands

	return result
}

func (c *Command) ParseGroups(initArgs, initConfigs []string) (results []*ParseResult) {
	leafCmd, commands, optionSetInitArgs, optionSetInitConfigs := c.splitCommandsArgs(initArgs, initConfigs)

	if len(optionSetInitArgs) == 0 && len(optionSetInitConfigs) == 0 {
		result := leafCmd.options.Parse(optionSetInitArgs, optionSetInitConfigs)
		results = append(results, result)
	} else {
		results = leafCmd.options.ParseGroups(optionSetInitArgs, optionSetInitConfigs)
	}

	for _, result := range results {
		result.commands = commands
	}

	return results
}

func (c *Command) GetHelp() []byte {
	buffer := &bytes.Buffer{}

	if len(c.name) > 0 {
		buffer.WriteString(path.Base(c.name))
		buffer.WriteString(": ")
	}
	if len(c.summary) > 0 {
		buffer.WriteString(c.summary)
	}
	if len(c.name) > 0 || len(c.summary) > 0 {
		buffer.WriteByte('\n')
	} else {
		buffer.WriteString("Usage:\n")
	}

	optionsHelp := c.options.GetHelp()
	if len(optionsHelp) > 0 {
		buffer.WriteString("\nOptions:\n\n")
		buffer.Write(optionsHelp)
	}

	if len(c.subCommands) > 0 {
		buffer.WriteString("\nSub commands:\n\n")
		for _, cmd := range c.subCommands {
			buffer.WriteString(cmd.name)
			buffer.WriteByte('\n')
			if len(cmd.summary) > 0 {
				buffer.WriteString(cmd.summary)
				buffer.WriteByte('\n')
			}
			buffer.WriteByte('\n')
		}
	}

	return buffer.Bytes()
}

func (c *Command) PrintHelp() {
	os.Stdout.Write(c.GetHelp())
}
