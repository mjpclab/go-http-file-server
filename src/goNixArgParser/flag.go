package goNixArgParser

func NewFlag(name string, canMerge, canFollowAssign, canConcatAssign bool, assignSigns []string) *Flag {
	return &Flag{
		Name:            name,
		canMerge:        canMerge,
		canFollowAssign: canFollowAssign,
		canConcatAssign: canConcatAssign,
		assignSigns:     assignSigns,
	}
}

func NewSimpleFlag(name string) *Flag {
	isSingleChar := len(name) == 1 || (len(name) == 2 && name[0] == '-')

	canMerge := isSingleChar
	canConcatAssign := isSingleChar

	assignSigns := make([]string, 0, 1)
	if !isSingleChar {
		assignSigns = append(assignSigns, "=")
	}

	return NewFlag(name, canMerge, true, canConcatAssign, assignSigns)
}

func NewSimpleFlags(names []string) []*Flag {
	flags := make([]*Flag, 0, len(names))

	for _, name := range names {
		flag := NewSimpleFlag(name)
		flags = append(flags, flag)
	}

	return flags
}
