package goVirtualHost

import "testing"

func TestCertsKeysToPairs(t *testing.T) {
	certFiles := []string{"cert1.crt", "cert2.crt", "cert3.crt"}
	keyFiles := []string{"cert1.key", "cert2.key", "cert3.key"}
	pairs, _ := CertsKeysToPairs(certFiles, keyFiles)

	if len(pairs) != 3 {
		t.Error(len(pairs))
	}

	if pairs[0][0] != "cert1.crt" ||
		pairs[0][1] != "cert1.key" ||
		pairs[1][0] != "cert2.crt" ||
		pairs[1][1] != "cert2.key" ||
		pairs[2][0] != "cert3.crt" ||
		pairs[2][1] != "cert3.key" {
		t.Error(pairs)
	}
}
