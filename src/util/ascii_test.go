package util

import "testing"

func TestAsciiToLowerCase(t *testing.T) {
	str := "Hello, 你好"
	lower := AsciiToLowerCase(str)
	expect := "hello, 你好"
	if lower != expect {
		t.Error(lower)
	}
}

func BenchmarkAsciiToLowerCase(b *testing.B) {
	const str = "Confucius said: \"三人行，必有我师焉\""
	for i := 0; i < b.N; i++ {
		AsciiToLowerCase(str)
	}
}
