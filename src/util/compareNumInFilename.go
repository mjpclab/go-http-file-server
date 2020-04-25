package util

import "bytes"

func findCommonPrefix(prev, next []byte) int {
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

func extractPrefixDigits(input []byte) []byte {
	buf := input[0:0]

	var prefixLen, length int
	for prefixLen, length = 0, len(input); prefixLen < length; prefixLen++ {
		b := input[prefixLen]
		if b < '0' || b > '9' {
			break
		}
	}
	return buf[:prefixLen]
}

func CompareNumInFilename(prev, next []byte) (less, ok bool) {
	if len(prev) == 0 && len(next) == 0 {
		return false, false
	} else if len(prev) == 0 {
		return true, true
	} else if len(next) == 0 {
		return false, true
	}

	common := findCommonPrefix(prev, next)
	if common > 0 {
		prev = prev[common:]
		next = next[common:]

		if len(prev) == 0 {
			return true, true
		} else if len(next) == 0 {
			return false, true
		}
	}

	prevDigits := extractPrefixDigits(prev)
	nextDigits := extractPrefixDigits(next)
	prevDigitsLen := len(prevDigits)
	nextDigitsLen := len(nextDigits)

	if prevDigitsLen != nextDigitsLen {
		return prevDigitsLen < nextDigitsLen, true
	}

	if prevDigitsLen == 0 { // prevDigitsLen and nextDigitsLen is 0
		switch {
		case prev[0] == '.' && next[0] != '.':
			return true, true
		case next[0] == '.' && prev[0] != '.':
			return false, true
		default:
			byteCmpResult := bytes.Compare(prev, next)
			if byteCmpResult != 0 {
				return byteCmpResult < 0, true
			} else {
				return false, false
			}
		}
	}

	compareResult := bytes.Compare(prevDigits, nextDigits)
	if compareResult != 0 {
		return compareResult < 0, true
	} else {
		return CompareNumInFilename(prev[prevDigitsLen:], next[nextDigitsLen:])
	}
}
