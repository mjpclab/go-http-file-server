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

func WildcardToRegexp(wildcard string) string {
	exp := "^" + regexpEscapeReplacer.Replace(wildcard) + "$"
	return exp
}
