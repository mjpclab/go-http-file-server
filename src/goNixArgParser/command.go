package goNixArgParser

import (
	"bytes"
	"os"
	"path"
)

func NewCommand(
	name, summary, mergeOptionPrefix string,
	restSigns []string,
) *Command {
	return &Command{
		Name:        name,
		Summary:     summary,
		OptionSet:   NewOptionSet(mergeOptionPrefix, restSigns),
		SubCommands: []*Command{},
	}
}

func NewSimpleCommand(name, summary string) *Command {
	return &Command{
		Name:        name,
		Summary:     summary,
		OptionSet:   NewSimpleOptionSet(),
		SubCommands: []*Command{},
	}
}

func (c *Command) NewSubCommand(
	name, summary, mergeOptionPrefix string,
	restSigns []string,
) *Command {
	subCommand := NewCommand(name, summary, mergeOptionPrefix, restSigns)
	c.SubCommands = append(c.SubCommands, subCommand)
	return subCommand
}

func (c *Command) NewSimpleSubCommand(name, summary string) *Command {
	subCommand := NewSimpleCommand(name, summary)
	c.SubCommands = append(c.SubCommands, subCommand)
	return subCommand
}

func (c *Command) GetSubCommand(name string) *Command {
	if c.SubCommands == nil {
		return nil
	}

	for _, cmd := range c.SubCommands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

func (c *Command) getNormalizedArgs(initArgs []string) ([]*Arg, *Command) {
	cmd := c

	if len(initArgs) == 0 {
		return []*Arg{}, cmd
	}

	args := make([]*Arg, 0, len(initArgs))

	for i, arg := range initArgs {
		if i == 0 && cmd.Name == arg {
			args = append(args, NewArg(arg, CommandArg))
		} else if subCmd := cmd.GetSubCommand(arg); subCmd != nil {
			args = append(args, NewArg(arg, CommandArg))
			cmd = subCmd
		} else {
			break
		}
	}

	return args, cmd
}

func (c *Command) Parse(initArgs []string) *ParseResult {
	args, cmd := c.getNormalizedArgs(initArgs)

	commands := []string{}
	for _, arg := range args {
		if arg.Type != CommandArg {
			break
		}
		commands = append(commands, arg.Text)
	}

	optionSetInitArgs := initArgs[len(args):]
	result := cmd.OptionSet.Parse(optionSetInitArgs)
	result.commands = commands

	return result
}

func (c *Command) GetHelp() []byte {
	buffer := &bytes.Buffer{}

	if len(c.Name) > 0 {
		buffer.WriteString(path.Base(c.Name))
		buffer.WriteString(": ")
	}
	if len(c.Summary) > 0 {
		buffer.WriteString(c.Summary)
	}
	if len(c.Name) > 0 || len(c.Summary) > 0 {
		buffer.WriteByte('\n')
	} else {
		buffer.WriteString("Usage:\n")
	}

	optionsHelp := c.OptionSet.GetHelp()
	if len(optionsHelp) > 0 {
		buffer.WriteString("\nOptions:\n\n")
		buffer.Write(optionsHelp)
	}

	if len(c.SubCommands) > 0 {
		buffer.WriteString("\nSub commands:\n\n")
		for _, cmd := range c.SubCommands {
			buffer.WriteString(cmd.Name)
			buffer.WriteByte('\n')
			if len(cmd.Summary) > 0 {
				buffer.WriteString(cmd.Summary)
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
