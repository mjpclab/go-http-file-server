package param

import (
	"../util"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func splitKeyValues(input string) (key string, values []string, ok bool) {
	sep, sepLen := utf8.DecodeRuneInString(input)
	if sepLen == 0 {
		return
	}
	entry := input[sepLen:]
	if len(entry) == 0 {
		return
	}

	sepIndex := strings.IndexRune(entry, sep)
	if sepIndex == 0 { // no key
		return
	} else if sepIndex > 0 {
		key = entry[:sepIndex]
		values = strings.FieldsFunc(entry[sepIndex+sepLen:], func(r rune) bool {
			return r == sep
		})
	} else { // only key
		key = entry
	}

	return key, values, true
}

func splitKeyValue(input string) (sep rune, sepLen int, k, v string, ok bool) {
	sep, sepLen = utf8.DecodeRuneInString(input)
	if sepLen == 0 {
		return
	}
	entry := input[sepLen:]
	if len(entry) == 0 {
		return
	}

	sepIndex := strings.IndexRune(entry, sep)
	if sepIndex <= 0 || sepIndex+sepLen == len(entry) {
		return
	}

	k = entry[:sepIndex]
	v = entry[sepIndex+sepLen:]
	return sep, sepLen, k, v, true
}

func splitAllKeyValue(inputs []string) (results [][2]string) {
	results = make([][2]string, 0, len(inputs))
	for i := range inputs {
		_, _, k, v, ok := splitKeyValue(inputs[i])
		if ok {
			results = append(results, [2]string{k, v})
		}
	}
	return
}

func normalizePathRestrictAccesses(
	inputs []string,
	normalizePath func(string) (string, error),
) (maps map[string][]string, errs []error) {
	maps = make(map[string][]string, len(inputs))

	for i := range inputs {
		reqPath, hosts, ok := splitKeyValues(inputs[i])
		if !ok {
			continue
		}

		normalizedPath, err := normalizePath(reqPath)
		if err != nil {
			errs = append(errs, err)
		}
		normalizedHosts := util.ExtractHostsFromUrls(hosts)

		for existingPath := range maps {
			if util.IsPathEqual(existingPath, normalizedPath) {
				normalizedPath = existingPath
				break
			}
		}

		maps[normalizedPath] = append(maps[normalizedPath], normalizedHosts...)
	}

	return
}

func normalizePathHeadersMap(
	inputs []string,
	normalizePath func(string) (string, error),
) (maps map[string][][2]string, errs []error) {
	maps = make(map[string][][2]string, len(inputs))

	for _, input := range inputs {
		sep, sepLen, reqPath, header, ok := splitKeyValue(input)
		if !ok {
			continue
		}
		sepIndex := strings.IndexRune(header, sep)
		if sepIndex <= 0 || sepIndex+sepLen == len(header) {
			continue
		}

		normalizedPath, err := normalizePath(reqPath)
		if err != nil {
			errs = append(errs, err)
		}
		headerName := header[:sepIndex]
		headerValue := header[sepIndex+1:]

		for existingPath := range maps {
			if util.IsPathEqual(existingPath, normalizedPath) {
				normalizedPath = existingPath
				break
			}
		}

		maps[normalizedPath] = append(maps[normalizedPath], [2]string{headerName, headerValue})
	}

	return
}

func normalizePathMaps(inputs [][2]string) (results [][2]string, errs []error) {
	results = make([][2]string, 0, len(inputs))

eachInput:
	for i := range inputs {
		urlPath := inputs[i][0]
		fsPath := inputs[i][1]
		if len(urlPath) == 0 || len(fsPath) == 0 {
			continue
		}
		urlPath = util.CleanUrlPath(urlPath)
		fsPath, err := util.NormalizeFsPath(fsPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for j := range results {
			if util.IsPathEqual(results[j][0], urlPath) {
				results[j][1] = fsPath
				continue eachInput
			}
		}

		results = append(results, [2]string{urlPath, fsPath})
	}

	return
}

func normalizeUrlPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}
		outputs = append(outputs, util.CleanUrlPath(input))
	}

	return outputs
}

func normalizeFsPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}

		abs, err := util.NormalizeFsPath(input)
		if err != nil {
			continue
		}

		outputs = append(outputs, abs)
	}

	return outputs
}

func normalizeFilenames(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}

		if strings.IndexByte(input, '/') >= 0 {
			continue
		}

		if filepath.Separator != '/' && strings.IndexByte(input, filepath.Separator) >= 0 {
			continue
		}

		outputs = append(outputs, input)
	}

	return outputs
}

func validateHstsPort(listensPlain, listensTLS []string) bool {
	var fromOK, toOK bool

	for _, listen := range listensPlain {
		port := util.ExtractListenPort(listen)
		if len(port) == 0 || port == "80" {
			fromOK = true
			break
		}
	}

	for _, listen := range listensTLS {
		port := util.ExtractListenPort(listen)
		if len(port) == 0 || port == "443" {
			toOK = true
			break
		}
	}

	return fromOK && toOK
}

func normalizeHttpsPort(httpsPort string, listensTLS []string) (string, bool) {
	if len(httpsPort) > 0 {
		httpsPort = util.ExtractListenPort(httpsPort)
		if len(httpsPort) == 0 {
			return "", false
		}
	} else if len(listensTLS) > 0 {
		httpsPort = util.ExtractListenPort(listensTLS[0])
	}

	for _, listen := range listensTLS {
		port := util.ExtractListenPort(listen)
		if len(httpsPort) == 0 && len(port) == 0 {
			return "", true
		}
		if httpsPort == port {
			return ":" + httpsPort, true
		}

		if httpsPort == "443" && port == "" {
			return "", true
		}
	}

	return "", false
}

func normalizeRedirectCode(code int) int {
	if code <= 300 || code > 399 {
		return 301
	}
	return code
}
