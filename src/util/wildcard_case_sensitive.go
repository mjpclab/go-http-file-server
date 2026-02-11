//go:build !windows && !darwin
// +build !windows,!darwin

package util

func WildcardToStrRegexp(wildcard string) string {
	exp := "^" + regexpEscapeReplacer.Replace(wildcard) + "$"
	return exp
}
