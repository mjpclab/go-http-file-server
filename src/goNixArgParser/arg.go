package goNixArgParser

func NewArg(text string, argType ArgType) *Arg {
	return &Arg{
		Text: text,
		Type: argType,
	}
}
