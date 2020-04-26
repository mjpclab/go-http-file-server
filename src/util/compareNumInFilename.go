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

func extractPrefixDigits(input []byte) ([]byte, int) {
	buf := input[0:0]

	var prefixLen, length int
	for prefixLen, length = 0, len(input); prefixLen < length; prefixLen++ {
		b := input[prefixLen]
		if b < '0' || b > '9' {
			break
		}
	}
	return buf[:prefixLen], prefixLen
}

func compareIgnoreAsciiCase(prev, next []byte) (less, ok bool) {
	prevLen := len(prev)
	nextLen := len(next)

	maxLen := prevLen
	if nextLen < maxLen {
		maxLen = nextLen
	}

	for i := 0; i < maxLen; i++ {
		prevByte := prev[i]
		prevChar := prevByte
		if prevChar >= 'A' && prevChar <= 'Z' {
			prevChar += 'a' - 'A'
		}

		nextByte := next[i]
		nextChar := nextByte
		if nextChar >= 'A' && nextChar <= 'Z' {
			nextChar += 'a' - 'A'
		}

		if prevChar != nextChar {
			return prevChar < nextChar, true
		} else if prevByte != nextByte {
			return prevByte < nextByte, true
		}
	}

	if prevLen != nextLen {
		return prevLen < nextLen, true
	}

	return
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

	prevDigits, prevDigitsLen := extractPrefixDigits(prev)
	nextDigits, nextDigitsLen := extractPrefixDigits(next)

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
			return compareIgnoreAsciiCase(prev, next)
		}
	}

	compareResult := bytes.Compare(prevDigits, nextDigits)
	if compareResult != 0 {
		return compareResult < 0, true
	} else {
		return CompareNumInFilename(prev[prevDigitsLen:], next[nextDigitsLen:])
	}
}
