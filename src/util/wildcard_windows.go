//go:build windows
// +build windows

package util

func WildcardToStrRegexp(wildcard string) string {
	exp := "(?i)^" + regexpEscapeReplacer.Replace(wildcard) + "$"
	return exp
}
