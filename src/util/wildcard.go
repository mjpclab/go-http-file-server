package util

import (
	"strings"
)

var regexpEscapeReplacer = strings.NewReplacer(
	"(", "\\(",
	")", "\\)",
	"[", "\\[",
	"]", "\\]",
	"{", "\\{",
	"}", "\\}",
	"<", "\\<",
	">", "\\>",
	"^", "\\^",
	"$", "\\$",
	"|", "\\|",
	"+", "\\+",
	"\\", "\\\\",
	".", "\\.",
	"?", ".",
	"*", ".*?",
)

func WildcardToStrRegexp(wildcard string) string {
	exp := "^" + regexpEscapeReplacer.Replace(wildcard) + "$"
	return exp
}
