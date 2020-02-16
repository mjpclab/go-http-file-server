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

func asciiToLowerCase(input string) string {
	buffer := []byte(input)
	length := len(buffer)

	for i := 0; i < length; {
		r, w := utf8.DecodeRune(buffer[i:])
		if w == 1 && r >= 'A' && r <= 'Z' {
			buffer[i] += 'a' - 'A'
		}

		i += w
	}

	return string(buffer)
}

func normalizeFsPath(input string) (string, error) {
	abs, err := filepath.Abs(input)
	if err != nil {
		return abs, err
	}

	volume := filepath.VolumeName(abs)
	if len(volume) > 0 {
		// suppose on windows platform, ignore ascii case in path name
		abs = asciiToLowerCase(abs)
	}

	return abs, err
}

func normalizeFsPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}

		abs, err := normalizeFsPath(input)
		if err != nil {
			continue
		}

		outputs = append(outputs, abs)
	}

	return outputs
}
