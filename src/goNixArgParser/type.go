package goNixArgParser

type Command struct {
	Name        string
	Summary     string
	OptionSet   *OptionSet
	SubCommands []*Command
}

type OptionSet struct {
	mergeFlagPrefix string
	restSigns       []string
	options         []*Option

	keyOptionMap  map[string]*Option
	flagOptionMap map[string]*Option
	flagMap       map[string]*Flag
	keyEnvMap     map[string][]string
	keyDefaultMap map[string][]string
}

type Option struct {
	Key           string
	Summary       string
	Description   string
	Flags         []*Flag
	AcceptValue   bool
	MultiValues   bool
	OverridePrev  bool
	Delimiters    []rune
	EnvVars       []string
	DefaultValues []string
}

type Flag struct {
	Name            string
	canMerge        bool
	canEqualAssign  bool
	canConcatAssign bool
}

type ArgType int

const (
	UnknownArg ArgType = iota
	CommandArg
	FlagArg
	ValueArg
	RestSignArg
	RestArg
)

type Arg struct {
	Text string
	Type ArgType
}

type ParseResult struct {
	commands []string
	params   map[string][]string
	envs     map[string][]string
	defaults map[string][]string
	rests    []string
}
