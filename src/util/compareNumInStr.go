package util

import "bytes"

func findCommonPrefix(prev, next string) int {
	prevLen := len(prev)
	nextLen := len(next)

	var maxLen int
	if prevLen < nextLen {
		maxLen = prevLen
	} else {
		maxLen = nextLen
	}

	for i := 0; i < maxLen; i++ {
		if prev[i] != next[i] {
			return i
		}
	}

	return maxLen
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func extractPrefixDigits(input string) (output string) {
	buf := bytes.Buffer{}
	for i, length := 0, len(input); i < length; i++ {
		b := input[i]
		if !isDigit(b) {
			break
		}
		buf.WriteByte(b)
	}
	return buf.String()
}

func CompareNumInStr(prev, next string) bool {
	if len(prev) == 0 {
		return true
	} else if len(next) == 0 {
		return false
	}

	common := findCommonPrefix(prev, next)
	if common > 0 {
		prev = prev[common:]
		next = next[common:]

		if len(prev) == 0 {
			return true
		} else if len(next) == 0 {
			return false
		}
	}

	prevDigits := extractPrefixDigits(prev)
	nextDigits := extractPrefixDigits(next)
	prevDigitsLen := len(prevDigits)
	nextDigitsLen := len(nextDigits)

	if prevDigitsLen != nextDigitsLen {
		return prevDigitsLen < nextDigitsLen
	}

	if prevDigitsLen == 0 {
		return prev < next
	}

	if prevDigits == nextDigits {
		return CompareNumInStr(prev[prevDigitsLen:], next[nextDigitsLen:])
	} else {
		return prevDigits < nextDigits
	}
}
