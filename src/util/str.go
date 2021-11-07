package util

import "strings"

func IsStrEqualAccurate(a, b string) bool {
	return a == b
}

func IsStrEqualNoCase(a, b string) bool {
	return strings.EqualFold(a, b)
}
