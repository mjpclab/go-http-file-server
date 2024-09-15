package util

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type StrEqualFunc func(a, b string) bool

func InPlaceDedup(inputs []string) []string {
	if len(inputs) <= 1 {
		return inputs
	}

	endIndex := 1
eachValue:
	for i, length := 1, len(inputs); i < length; i++ {
		for j := 0; j < endIndex; j++ {
			if inputs[i] == inputs[j] {
				continue eachValue
			}
		}
		if endIndex != i {
			inputs[endIndex] = inputs[i]
		}
		endIndex++
	}

	return inputs[:endIndex]
}

func InPlaceTrim(inputs []string) {
	for i := range inputs {
		inputs[i] = strings.TrimSpace(inputs[i])
	}
}

func EscapeControllingRune(str string) []byte {
	runeBytes := make([]byte, 4)
	buf := make([]byte, 0, len(str))

	for _, r := range str {
		if uint32(r) < 32 { // non-printable chars
			b := byte(r)
			if b == 0 {
				buf = append(buf, '\\', '0')
			} else if b == '\a' {
				buf = append(buf, '\\', 'a')
			} else if b == '\b' {
				buf = append(buf, '\\', 'b')
			} else if b == '\f' {
				buf = append(buf, '\\', 'f')
			} else if b == '\n' {
				buf = append(buf, '\\', 'n')
			} else if b == '\r' {
				buf = append(buf, '\\', 'r')
			} else if b == '\t' {
				buf = append(buf, '\\', 't')
			} else if b == '\v' {
				buf = append(buf, '\\', 'v')
			} else {
				h, l := ByteToHex(b)
				buf = append(buf, '\\', 'x', h, l)
			}
			continue
		}

		if unicode.IsControl(r) {
			nBytes := utf8.EncodeRune(runeBytes, r)
			for i := 0; i < nBytes; i++ {
				h, l := ByteToHex(runeBytes[i])
				buf = append(buf, '\\', 'x', h, l)
			}
		} else {
			buf = utf8.AppendRune(buf, r)
		}
	}

	return buf
}
