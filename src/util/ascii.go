package util

import "unicode/utf8"

func AsciiToLowerCase(input string) string {
	buffer := []byte(input)
	length := len(buffer)

	for i := 0; i < length; {
		r, w := utf8.DecodeRune(buffer[i:])
		if w == 1 && r >= 'A' && r <= 'Z' {
			buffer[i] += 'a' - 'A'
		}

		i += w
	}

	return string(buffer)
}
