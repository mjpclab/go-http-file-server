package goNixArgParser

func newToken(text string, argType argKind) *argToken {
	return &argToken{
		text: text,
		kind: argType,
	}
}
