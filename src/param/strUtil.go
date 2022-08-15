package param

import (
	"mjpclab.dev/ghfs/src/util"
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

func splitAllKeyValues(inputs []string) (results [][]string) {
	results = make([][]string, 0, len(inputs))
	for i := range inputs {
		key, values, ok := splitKeyValues(inputs[i])
		if !ok {
			continue
		}
		keyValues := make([]string, 1+len(values))
		keyValues[0] = key
		copy(keyValues[1:], values)
		results = append(results, keyValues)
	}
	return
}

func splitKeyValue(input string) (k, v string, ok bool) {
	sep, sepLen := utf8.DecodeRuneInString(input)
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
	return k, v, true
}

func splitAllKeyValue(inputs []string) (results [][2]string) {
	results = make([][2]string, 0, len(inputs))
	for i := range inputs {
		k, v, ok := splitKeyValue(inputs[i])
		if ok {
			results = append(results, [2]string{k, v})
		}
	}
	return
}

func normalizeAllPathValues(
	inputs [][]string,
	keepEmptyValuesEntry bool,
	normalizePath func(string) (string, error),
	normalizeEntryValues func([]string) []string,
) (results [][]string, errs []error) {
	var err error
	results = make([][]string, 0, len(inputs))

eachInput:
	for i := range inputs {
		if len(inputs[i]) == 0 {
			continue
		}

		reqPath := inputs[i][0]
		if len(reqPath) == 0 {
			continue
		}
		reqPath, err = normalizePath(reqPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		values := inputs[i][1:]
		if normalizeEntryValues != nil {
			values = normalizeEntryValues(values)
		}
		for j := range results {
			if util.IsPathEqual(results[j][0], reqPath) {
				results[j] = append(results[j], values...)
				continue eachInput
			}
		}

		if len(values) == 0 && !keepEmptyValuesEntry {
			continue
		}

		entry := make([]string, 1+len(values))
		entry[0] = reqPath
		copy(entry[1:], values)

		results = append(results, entry)
	}

	return
}

func dedupPathValues(inputs []string) []string {
	if len(inputs) <= 2 { // path & single value
		return inputs
	}

	values := inputs[1:]
	endIndex := 1
eachValue:
	for i, length := 1, len(values); i < length; i++ {
		for j := 0; j < endIndex; j++ {
			if values[i] == values[j] {
				continue eachValue
			}
		}
		if endIndex != i {
			values[endIndex] = values[i]
		}
		endIndex++
	}

	return inputs[:1+endIndex]
}

func dedupAllPathValues(inputs [][]string) {
	for i, iLen := 0, len(inputs); i < iLen; i++ {
		inputs[i] = dedupPathValues(inputs[i])
	}
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

func normalizeHeaders(inputs []string) []string {
	if len(inputs) != 2 {
		return nil
	}
	if len(inputs[0]) == 0 || len(inputs[1]) == 0 {
		return nil
	}
	return inputs
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
