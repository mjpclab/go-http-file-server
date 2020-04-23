package goNixArgParser

type Command struct {
	names       []string
	summary     string
	options     *OptionSet
	subCommands []*Command
}

type OptionSet struct {
	mergeFlagPrefix   string
	restsSigns        []string
	groupSeps         []string
	undefFlagPrefixes []string

	options []*Option

	hasCanMerge        bool
	hasCanConcatAssign bool
	hasAssignSigns     bool
	hasPrefixMatch     bool

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
	prefixMatchLen  int
	canMerge        bool
	canFollowAssign bool
	canConcatAssign bool
	assignSigns     []string
}

type ArgType int

const (
	UndetermArg ArgType = iota
	CommandArg
	FlagArg
	ValueArg
	UndefFlagArg
	UndefFlagValueArg
	AmbiguousFlagArg
	AmbiguousFlagValueArg
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

	argAmbigus    []string
	configAmbigus []string

	argUndefs    []string
	configUndefs []string
}
