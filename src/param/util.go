package param

import (
	"mjpclab.dev/ghfs/src/util"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// SplitKeyValues
// input format: <sep><key>[<sep><value>...]
func SplitKeyValues(input string) (key string, values []string, ok bool) {
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

func SplitAllKeyValues(inputs []string) (results [][]string) {
	results = make([][]string, 0, len(inputs))
	for i := range inputs {
		key, values, ok := SplitKeyValues(inputs[i])
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

// SplitKeyValue
// input format: <sep><key><sep><value>
func SplitKeyValue(input string) (k, v string, ok bool) {
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

func SplitAllKeyValue(inputs []string) (results [][2]string) {
	results = make([][2]string, 0, len(inputs))
	for i := range inputs {
		k, v, ok := SplitKeyValue(inputs[i])
		if ok {
			results = append(results, [2]string{k, v})
		}
	}
	return
}

// EntriesToKVs
// input element: "key:value"
// output element: [2]string{"key, "value"}
func EntriesToKVs(entries []string) [][2]string {
	KVs := make([][2]string, 0, len(entries))
	for _, entry := range entries {
		colonIndex := strings.IndexByte(entry, ':')
		if colonIndex <= 0 || colonIndex == len(entry)-1 {
			continue
		}
		key := entry[:colonIndex]
		value := entry[colonIndex+1:]
		KVs = append(KVs, [2]string{key, value})
	}
	return KVs
}

func NormalizeUrlPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}
		outputs = append(outputs, util.CleanUrlPath(input))
	}

	outputs = util.InPlaceDedup(outputs)

	return outputs
}

func NormalizeFsPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}

		abs, err := filepath.Abs(input)
		if err != nil {
			continue
		}

		outputs = append(outputs, abs)
	}

	outputs = util.InPlaceDedup(outputs)

	return outputs
}

func NormalizeRedirectCode(code int) int {
	if code <= 300 || code > 399 {
		return 301
	}
	return code
}
