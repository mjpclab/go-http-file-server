package shimgo

import (
	"crypto/tls"
	"encoding/base64"
	"net"
	"net/http"
	"strings"
)

const (
	Net_Http_MethodGet     = "GET"
	Net_Http_MethodHead    = "HEAD"
	Net_Http_MethodPost    = "POST"
	Net_Http_MethodPut     = "PUT"
	Net_Http_MethodPatch   = "PATCH" // RFC 5789
	Net_Http_MethodDelete  = "DELETE"
	Net_Http_MethodConnect = "CONNECT"
	Net_Http_MethodOptions = "OPTIONS"
	Net_Http_MethodTrace   = "TRACE"
)

func Net_Http_Server_ServeTLS(srv *http.Server, l net.Listener, certFile, keyFile string) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	if len(config.Certificates) == 0 {
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
	}

	tlsListener := tls.NewListener(l, config)
	return srv.Serve(tlsListener)
}

func Net_Http_BasicAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}
	return net_http_parseBasicAuth(auth)
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func net_http_parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
