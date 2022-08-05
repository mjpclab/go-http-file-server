package param

import (
	"../goVirtualHost"
	"crypto/tls"
	"strings"
)

func loadCertificates(certFiles, keyFiles []string) ([]tls.Certificate, []error) {
	return goVirtualHost.LoadCertificates(certFiles, keyFiles)
}

func entriesToUsers(entries []string) [][2]string {
	users := make([][2]string, 0, len(entries))
	for _, userEntry := range entries {
		username := userEntry
		password := ""

		colonIndex := strings.IndexByte(userEntry, ':')
		if colonIndex >= 0 {
			username = userEntry[:colonIndex]
			password = userEntry[colonIndex+1:]
		}

		users = append(users, [2]string{username, password})
	}
	return users
}

func entriesToHeaders(entries []string) [][2]string {
	headers := make([][2]string, 0, len(entries))
	for _, entry := range entries {
		colonIndex := strings.IndexByte(entry, ':')
		if colonIndex <= 0 || colonIndex == len(entry)-1 {
			continue
		}
		key := entry[:colonIndex]
		value := entry[colonIndex+1:]
		headers = append(headers, [2]string{key, value})
	}
	return headers
}
