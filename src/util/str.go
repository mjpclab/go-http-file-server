package util

import "strings"

type StrEqualFunc func(a, b string) bool

func IsStrEqualAccurate(a, b string) bool {
	return a == b
}

func IsStrEqualNoCase(a, b string) bool {
	return strings.EqualFold(a, b)
}
