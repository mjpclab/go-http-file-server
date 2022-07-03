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
	assignSigns       []string
	undefFlagPrefixes []string

	options []*Option

	hasCanMerge        bool
	hasCanConcatAssign bool
	hasPrefixMatch     bool

	keyOptionMap  map[string]*Option
	flagOptionMap map[string]*Option
	nameFlagMap   map[string]*Flag
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
	Hidden        bool
}

type Flag struct {
	Name            string
	prefixMatchLen  int
	canMerge        bool
	canFollowAssign bool
	canConcatAssign bool
}

type argKind int

const (
	undetermArg argKind = iota
	commandArg
	flagArg
	valueArg
	undefFlagArg
	undefFlagValueArg
	ambiguousFlagArg
	ambiguousFlagValueArg
	restSignArg
	restArg
	groupSepArg
)

type argToken struct {
	text string
	kind argKind
}

type ParseResult struct {
	keyOptionMap map[string]*Option

	commands         []string
	specifiedOptions map[string][]string
	envs             map[string][]string
	configOptions    map[string][]string
	defaults         map[string][]string

	specifiedRests []string
	configRests    []string

	specifiedAmbigus []string
	configAmbigus    []string

	specifiedUndefs []string
	configUndefs    []string
}
