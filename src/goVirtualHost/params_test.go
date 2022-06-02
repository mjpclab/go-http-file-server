package goVirtualHost

import (
	"crypto/tls"
	"errors"
	"testing"
)

func TestParamsValidateParam(t *testing.T) {
	var p *param
	var ps params
	var errs []error

	// normal wildcard ip
	p = &param{
		proto: "tcp",
		ip:    "",
		port:  ":80",
	}
	errs = ps.validateParam(p)
	if len(errs) > 0 {
		t.Error(errs)
	}
	ps = append(ps, p)

	// same wildcard ip:port, different hostname
	p = &param{
		proto:     "tcp",
		ip:        "",
		port:      ":80",
		hostNames: []string{"localhost"},
	}
	errs = ps.validateParam(p)
	if len(errs) > 0 {
		t.Error(errs)
	}
	ps = append(ps, p)

	// IPv4 wildcard 0.0.0.0:port, conflict
	p = &param{
		proto: "tcp4",
		ip:    "",
		port:  ":80",
	}
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	} else if !errors.Is(errs[0], ConflictIPAddress) {
		t.Error()
	}

	// IPv6 wildcard [::]:port, conflict
	p = &param{
		proto: "tcp6",
		ip:    "",
		port:  ":80",
	}
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	} else if !errors.Is(errs[0], ConflictIPAddress) {
		t.Error()
	}

	// duplicated address and hostname
	p = &param{
		proto:     "tcp",
		ip:        "",
		port:      ":80",
		hostNames: []string{"localhost"},
	}
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error()
	} else if !errors.Is(errs[0], DuplicatedAddressHostname) {
		t.Error()
	}

	// cannot serve for both Plain and TLS mode
	p = &param{
		proto:  "tcp",
		ip:     "",
		port:   ":80",
		useTLS: true,
		certs:  []tls.Certificate{},
	}
	errs = ps.validateParam(p)
	if len(errs) == 0 {
		t.Error(errs)
	} else if !errors.Is(errs[0], ConflictTLSMode) {
		t.Error()
	}
}
