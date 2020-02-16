package param

import "testing"

func TestAsciiToLowerCase(t *testing.T) {
	str := "Hello, 你好"
	lower := asciiToLowerCase(str)
	expect := "hello, 你好"
	if lower != expect {
		t.Error(lower)
	}
}
