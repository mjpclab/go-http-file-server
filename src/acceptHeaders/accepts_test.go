package acceptHeaders

import "testing"

func TestParseAccepts(t *testing.T) {
	acceptLanguage := "zh;q=0.9,zh-CN,en;q=0.7,en-US;q=0.8"
	accepts := ParseAccepts(acceptLanguage)
	if len(accepts) != 4 {
		t.Error(len(accepts))
	}

	if accepts[0].value != "zh-CN" {
		t.Error(accepts[0].value)
	}
	if accepts[0].quality != 1000 {
		t.Error(accepts[0].quality)
	}

	if accepts[1].value != "zh" {
		t.Error(accepts[1].value)
	}
	if accepts[1].quality != 900 {
		t.Error(accepts[1].quality)
	}

	if accepts[2].value != "en-US" {
		t.Error(accepts[2].value)
	}
	if accepts[2].quality != 800 {
		t.Error(accepts[2].quality)
	}

	if accepts[3].value != "en" {
		t.Error(accepts[3].value)
	}
	if accepts[3].quality != 700 {
		t.Error(accepts[3].quality)
	}
}

func TestParseAccepts2(t *testing.T) {
	acceptEncoding := "gzip, deflate"
	accepts := ParseAccepts(acceptEncoding)

	if len(accepts) != 2 {
		t.Error(len(accepts))
	}

	if accepts[0].value != "gzip" {
		t.Error(accepts[0].value)
	}
	if accepts[0].quality != 1000 {
		t.Error(accepts[0].quality)
	}

	if accepts[1].value != "deflate" {
		t.Error(accepts[1].value)
	}
	if accepts[1].quality != 1000 {
		t.Error(accepts[1].quality)
	}
}

func TestGetPreferredValue(t *testing.T) {
	acceptEncoding := "gzip;q=0.9, deflate"
	accepts := ParseAccepts(acceptEncoding)

	var index int
	var preferred string
	index, preferred, _ = accepts.GetPreferredValue([]string{"gzip"})
	if index != 0 {
		t.Error(index)
	}
	if preferred != "gzip" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"deflate"})
	if index != 0 {
		t.Error(index)
	}
	if preferred != "deflate" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"deflate", "gzip"})
	if index != 0 {
		t.Error(index)
	}
	if preferred != "deflate" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"gzip", "deflate"})
	if index != 1 {
		t.Error(index)
	}
	if preferred != "deflate" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"gzip", "xxx"})
	if index != 0 {
		t.Error(index)
	}
	if preferred != "gzip" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"deflate", "xxx"})
	if index != 0 {
		t.Error(index)
	}
	if preferred != "deflate" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"xxx", "gzip", "deflate"})
	if index != 2 {
		t.Error(index)
	}
	if preferred != "deflate" {
		t.Error(preferred)
	}
}
