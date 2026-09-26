package serverHandler

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestExtractNoAuthUrl(t *testing.T) {
	h := &aliasHandler{}
	const fallback = "/dir/"

	cases := []struct {
		authValue string
		referer   string
		expect    string
	}{
		{"/dir/file", "", "/dir/file"},
		{"/dir/?sort=n", "", "/dir/?sort=n"},
		{"http://example.com/x/?a=1", "", "/x/?a=1"},
		{"relative/path", "", "/dir/relative/path"},
		{"../up/", "", "/up/"},
		{"/a b/中文/", "", "/a%20b/%E4%B8%AD%E6%96%87/"},
		{"/\\evil.example", "", "/%5Cevil.example"},
		{"", "http://example.com/other/?sort=s", "/other/?sort=s"},
		{"", "/other/?sort=s", "/other/?sort=s"},
		{"", "other/", "/dir/other/"},
		{"", "../x/?s=1", "/x/?s=1"},
		{"", "?sort=s", "/dir/?sort=s"},
		{"", "", fallback},

		{"https://evil.example", "", fallback},
		{"//evil.example", "", fallback},
		{".//evil.example/", "", "/dir//evil.example/"},
		{"/\t/evil.example", "", fallback},
		{"/\t\r\n/evil.example", "", fallback},
		{"/dir\twith\x7fctrl/", "", fallback},
		{"", "https://evil.example/", fallback},
		{"", "http://example.com//evil.example/", fallback},
		{"", "..//evil.example/", "/evil.example/"},
		{"", "mailto:a@example.com", fallback},
		{"https://evil.example", "https://evil.example/", fallback},
	}

	for _, c := range cases {
		r := httptest.NewRequest("GET", "http://example.com/dir/", nil)
		if len(c.referer) > 0 {
			r.Header.Set("Referer", c.referer)
		}
		session := &sessionContext{
			prefixReqPath: "/dir/",
			query:         url.Values{authQueryParam: {c.authValue}},
		}
		data := &responseData{}

		returnUrl := h.extractNoAuthUrl(r, session, data)
		if returnUrl != c.expect {
			t.Errorf("auth=%q referer=%q: got %q, expect %q", c.authValue, c.referer, returnUrl, c.expect)
		}
	}
}
