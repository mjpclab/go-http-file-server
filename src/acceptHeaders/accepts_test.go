package acceptHeaders

import (
	"reflect"
	"testing"
)

func TestParseAccept(t *testing.T) {
	accept := "text/html, */*;q=0.6, text/plain; q=0.9, application/json;q=0.7, image/*;q=0.8, image/png;q=0.8"
	accepts := ParseAccepts(accept)
	if len(accepts) != 6 {
		t.Error(len(accepts))
	}

	if !reflect.DeepEqual(accepts[0], acceptItem{"text/html", 1000, 0}) {
		t.Error(accepts[0])
	}
	if !reflect.DeepEqual(accepts[1], acceptItem{"text/plain", 900, 0}) {
		t.Error(accepts[1])
	}
	if !reflect.DeepEqual(accepts[2], acceptItem{"image/png", 800, 0}) {
		t.Error(accepts[2])
	}
	if !reflect.DeepEqual(accepts[3], acceptItem{"image/*", 800, 1}) {
		t.Error(accepts[3])
	}
	if !reflect.DeepEqual(accepts[4], acceptItem{"application/json", 700, 0}) {
		t.Error(accepts[4])
	}
	if !reflect.DeepEqual(accepts[5], acceptItem{"*/*", 600, 2}) {
		t.Error(accepts[5])
	}

	var index int
	var preferred string

	index, preferred, _ = accepts.GetPreferredValue([]string{"text/plain", "text/html"})
	if index != 1 {
		t.Error(index)
	}
	if preferred != "text/html" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"image/jpeg", "image/png"})
	if index != 1 {
		t.Error(index)
	}
	if preferred != "image/png" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"image/png", "image/jpeg"})
	if index != 0 {
		t.Error(index)
	}
	if preferred != "image/png" {
		t.Error(preferred)
	}

	index, preferred, _ = accepts.GetPreferredValue([]string{"image/webp", "image/jpeg"})
	if index != 0 {
		t.Error(index)
	}
	if preferred != "image/webp" {
		t.Error(preferred)
	}
}

func TestParseAcceptLanguage(t *testing.T) {
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

func TestParseAcceptEncoding(t *testing.T) {
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

func TestParseAcceptEncoding2(t *testing.T) {
	acceptEncoding := "gzip;v=b3;q=0.9, deflate"
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
