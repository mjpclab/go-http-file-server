package util

func AsciiToLowerCase(input string) string {
	buffer := []byte(input)

	for i := range buffer {
		if buffer[i] >= 'A' && buffer[i] <= 'Z' {
			buffer[i] += 'a' - 'A'
		}
	}

	return string(buffer)
}
