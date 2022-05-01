package util

const upperLowerDistance byte = 'a' - 'A'

func AsciiToLowerCase(input string) string {
	buffer := []byte(input)

	for i := range buffer {
		if buffer[i] >= 'A' && buffer[i] <= 'Z' {
			buffer[i] += upperLowerDistance
		}
	}

	return string(buffer)
}
