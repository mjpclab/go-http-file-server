//go:build windows || darwin
// +build windows darwin

package util

func WildcardToStrRegexp(wildcard string) string {
	exp := "(?i)^" + regexpEscapeReplacer.Replace(wildcard) + "$"
	return exp
}
