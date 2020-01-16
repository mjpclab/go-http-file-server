package param

import (
	"../util"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func splitMapping(input string) (k, v string, ok bool) {
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

func normalizePathMaps(inputs []string) map[string]string {
	maps := map[string]string{}

	for _, input := range inputs {
		urlPath, fsPath, ok := splitMapping(input)
		if !ok {
			continue
		}

		urlPath = util.CleanUrlPath(urlPath)
		fsPath = path.Clean(fsPath)
		maps[urlPath] = fsPath
	}

	return maps
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

		abs, err := filepath.Abs(input)
		if err != nil {
			continue
		}

		outputs = append(outputs, abs)
	}

	return outputs
}
