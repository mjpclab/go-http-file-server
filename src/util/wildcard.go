package util

import (
	"strings"
)

func WildcardToRegexp(wildcard string) string {
	replacer := strings.NewReplacer(
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

	exp := "^" + replacer.Replace(wildcard) + "$"

	return exp
}
