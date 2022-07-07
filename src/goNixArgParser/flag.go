package goNixArgParser

func NewFlag(name string, prefixMatchLen int, canMerge, canFollowAssign, canConcatAssign bool) *Flag {
	return &Flag{
		Name:            name,
		prefixMatchLen:  prefixMatchLen,
		canMerge:        canMerge,
		canFollowAssign: canFollowAssign,
		canConcatAssign: canConcatAssign,
	}
}

func NewSimpleFlag(name string) *Flag {
	isSingleChar := len(name) == 1 || (len(name) == 2 && name[0] == '-')

	canMerge := isSingleChar
	canConcatAssign := isSingleChar

	return NewFlag(name, 0, canMerge, true, canConcatAssign)
}

func NewSimpleFlags(names []string) []*Flag {
	flags := make([]*Flag, 0, len(names))

	for _, name := range names {
		flag := NewSimpleFlag(name)
		flags = append(flags, flag)
	}

	return flags
}
