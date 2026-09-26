package goVirtualHost

import (
	"crypto/tls"
	"crypto/x509"
	"testing"
)

func TestCertsMakeLeaf(t *testing.T) {
	certNoLeaf := newTestCert(t, "a.example.com")

	certWithLeaf := newTestCert(t, "b.example.com")
	existingLeaf := &x509.Certificate{}
	certWithLeaf.Leaf = existingLeaf

	certEmpty := &tls.Certificate{}
	certInvalid := &tls.Certificate{Certificate: [][]byte{[]byte("invalid")}}

	certs{certNoLeaf, certWithLeaf, certEmpty, certInvalid}.makeLeaf()

	if certNoLeaf.Leaf == nil || certNoLeaf.Leaf.Subject.CommonName != "a.example.com" {
		t.Error("leaf should be parsed", certNoLeaf.Leaf)
	}
	if certWithLeaf.Leaf != existingLeaf {
		t.Error("existing leaf should be kept")
	}
	if certEmpty.Leaf != nil {
		t.Error("certificate without data should have no leaf")
	}
	if certInvalid.Leaf != nil {
		t.Error("invalid certificate should have no leaf")
	}
}

func TestCertsMakeLeafEmpty(t *testing.T) {
	certs(nil).makeLeaf()
	certs{}.makeLeaf()
}
