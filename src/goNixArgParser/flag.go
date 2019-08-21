package goNixArgParser

func NewFlag(name string, canMerge, canEqualAssign, canConcatAssign bool) *Flag {
	return &Flag{
		Name:            name,
		canMerge:        canMerge,
		canEqualAssign:  canEqualAssign,
		canConcatAssign: canConcatAssign,
	}
}

func NewSimpleFlag(name string) *Flag {
	isSingleChar := len(name) == 1 || (len(name) == 2 && name[0] == '-')
	return NewFlag(name, isSingleChar, !isSingleChar, false)
}

func NewSimpleFlags(names []string) []*Flag {
	flags := make([]*Flag, 0, len(names))

	for _, name := range names {
		flag := NewSimpleFlag(name)
		flags = append(flags, flag)
	}

	return flags
}
