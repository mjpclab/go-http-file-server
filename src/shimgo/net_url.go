package shimgo

import (
	"bytes"
	"strconv"
)

type Net_Url_EscapeError string

func (e Net_Url_EscapeError) Error() string {
	return "invalid URL escape " + strconv.Quote(string(e))
}

type Net_Url_InvalidHostError string

func (e Net_Url_InvalidHostError) Error() string {
	return "invalid character " + strconv.Quote(string(e)) + " in host name"
}

type net_url_encoding int

const (
	net_url_encodePath net_url_encoding = 1 + iota
	net_url_encodePathSegment
	net_url_encodeHost
	net_url_encodeZone
	net_url_encodeUserPassword
	net_url_encodeQueryComponent
	net_url_encodeFragment
)
const net_url_upperhex = "0123456789ABCDEF"

func net_url_ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func net_url_unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func Net_Url_PathEscape(s string) string {
	return net_url_escape(s, net_url_encodePathSegment)
}

func net_url_escape(s string, mode net_url_encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if net_url_shouldEscape(c, mode) {
			if c == ' ' && mode == net_url_encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	if hexCount == 0 {
		copy(t, s)
		for i := 0; i < len(s); i++ {
			if s[i] == ' ' {
				t[i] = '+'
			}
		}
		return string(t)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == net_url_encodeQueryComponent:
			t[j] = '+'
			j++
		case net_url_shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = net_url_upperhex[c>>4]
			t[j+2] = net_url_upperhex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
//
// Please be informed that for now shouldEscape does not check all
// reserved characters correctly. See golang.org/issue/5684.
func net_url_shouldEscape(c byte, mode net_url_encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == net_url_encodeHost || mode == net_url_encodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case net_url_encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last three as well. That leaves only ? to escape.
			return c == '?'

		case net_url_encodePathSegment: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments.
			return c == '/' || c == ';' || c == ',' || c == '?'

		case net_url_encodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case net_url_encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case net_url_encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	if mode == net_url_encodeFragment {
		// RFC 3986 §2.2 allows not escaping sub-delims. A subset of sub-delims are
		// included in reserved from RFC 2396 §2.2. The remaining sub-delims do not
		// need to be escaped. To minimize potential breakage, we apply two restrictions:
		// (1) we always escape sub-delims outside of the fragment, and (2) we always
		// escape single quote to avoid breaking callers that had previously assumed that
		// single quotes would be escaped. See issue #19917.
		switch c {
		case '!', '(', ')', '*':
			return false
		}
	}

	// Everything else must be escaped.
	return true
}

// PathUnescape does the inverse transformation of PathEscape,
// converting each 3-byte encoded substring of the form "%AB" into the
// hex-decoded byte 0xAB. It returns an error if any % is not followed
// by two hexadecimal digits.
//
// PathUnescape is identical to QueryUnescape except that it does not
// unescape '+' to ' ' (space).
func Net_Url_PathUnescape(s string) (string, error) {
	return unescape(s, net_url_encodePathSegment)
}

// unescape unescapes a string; the mode specifies
// which section of the URL string is being unescaped.
func unescape(s string, mode net_url_encoding) (string, error) {
	// Count %, check that they're well-formed.
	n := 0
	hasPlus := false
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++
			if i+2 >= len(s) || !net_url_ishex(s[i+1]) || !net_url_ishex(s[i+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[:3]
				}
				return "", Net_Url_EscapeError(s)
			}
			// Per https://tools.ietf.org/html/rfc3986#page-21
			// in the host component %-encoding can only be used
			// for non-ASCII bytes.
			// But https://tools.ietf.org/html/rfc6874#section-2
			// introduces %25 being allowed to escape a percent sign
			// in IPv6 scoped-address literals. Yay.
			if mode == net_url_encodeHost && net_url_unhex(s[i+1]) < 8 && s[i:i+3] != "%25" {
				return "", Net_Url_EscapeError(s[i : i+3])
			}
			if mode == net_url_encodeZone {
				// RFC 6874 says basically "anything goes" for zone identifiers
				// and that even non-ASCII can be redundantly escaped,
				// but it seems prudent to restrict %-escaped bytes here to those
				// that are valid host name bytes in their unescaped form.
				// That is, you can use escaping in the zone identifier but not
				// to introduce bytes you couldn't just write directly.
				// But Windows puts spaces here! Yay.
				v := net_url_unhex(s[i+1])<<4 | net_url_unhex(s[i+2])
				if s[i:i+3] != "%25" && v != ' ' && net_url_shouldEscape(v, net_url_encodeHost) {
					return "", Net_Url_EscapeError(s[i : i+3])
				}
			}
			i += 3
		case '+':
			hasPlus = mode == net_url_encodeQueryComponent
			i++
		default:
			if (mode == net_url_encodeHost || mode == net_url_encodeZone) && s[i] < 0x80 && net_url_shouldEscape(s[i], mode) {
				return "", Net_Url_InvalidHostError(s[i : i+1])
			}
			i++
		}
	}

	if n == 0 && !hasPlus {
		return s, nil
	}

	var t bytes.Buffer
	t.Grow(len(s) - 2*n)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '%':
			t.WriteByte(net_url_unhex(s[i+1])<<4 | net_url_unhex(s[i+2]))
			i += 2
		case '+':
			if mode == net_url_encodeQueryComponent {
				t.WriteByte(' ')
			} else {
				t.WriteByte('+')
			}
		default:
			t.WriteByte(s[i])
		}
	}
	return t.String(), nil
}