package goNixArgParser

type OptionSet struct {
	mergeOptionPrefix string
	options           []*Option
	keyOptionMap      map[string]*Option
	flagOptionMap     map[string]*Option
	flagMap           map[string]*Flag
	keyDefaultMap     map[string][]string
}

type Option struct {
	Key          string
	Summary      string
	Description  string
	Flags        []*Flag
	AcceptValue  bool
	MultiValues  bool
	OverridePrev bool
	Delimiter    string
	DefaultValue []string
}
type Flag struct {
	Name            string
	canMerge        bool
	canEqualAssign  bool
	canConcatAssign bool
}

type ParseResult struct {
	params   map[string][]string
	defaults map[string][]string
	rests    []string
}

type ArgType int

const UnknownArg ArgType = 0
const FlagArg ArgType = 1
const ValueArg ArgType = 2
const RestArg ArgType = 3

type Arg struct {
	Text string
	Type ArgType
}
