package util

func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func IsDigits(input string) bool {
	for i, length := 0, len(input); i < length; i++ {
		b := input[i]
		if b < '0' || b > '9' {
			return false
		}
	}

	return true
}
