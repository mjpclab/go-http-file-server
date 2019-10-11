package goNixArgParser

type Command struct {
	name        string
	summary     string
	options     *OptionSet
	subCommands []*Command
}

type OptionSet struct {
	mergeFlagPrefix string
	restsSigns      []string
	groupSeps       []string

	options []*Option

	hasCanMerge        bool
	hasCanEqualAssign  bool
	hasCanConcatAssign bool

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
	UniqueValues  bool
	EnvVars       []string
	DefaultValues []string
}

type Flag struct {
	Name            string
	canMerge        bool
	canFollowAssign bool
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
	GroupSepArg
)

type Arg struct {
	Text string
	Type ArgType
}

type ParseResult struct {
	keyOptionMap map[string]*Option

	commands []string
	args     map[string][]string
	envs     map[string][]string
	configs  map[string][]string
	defaults map[string][]string

	argRests    []string
	configRests []string
}
