package param

import (
	"../util"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Param struct {
	Root          string
	Aliases       map[string]string
	GlobalUpload  bool
	Uploads       []string
	GlobalArchive bool
	Archives      []string
	Key           string
	Cert          string
	Listen        string
	Template      string
	Shows         *regexp.Regexp
	ShowDirs      *regexp.Regexp
	ShowFiles     *regexp.Regexp
	Hides         *regexp.Regexp
	HideDirs      *regexp.Regexp
	HideFiles     *regexp.Regexp
	AccessLog     string
	ErrorLog      string
}

func normalizePathMaps(inputs []string) map[string]string {
	maps := map[string]string{}

	for _, input := range inputs {
		sep, sepLen := utf8.DecodeRuneInString(input)
		if sepLen == 0 {
			continue
		}
		input = input[sepLen:]
		if len(input) == 0 {
			continue
		}

		sepIndex := strings.Index(input, string(sep))
		if sepIndex == -1 {
			continue
		}

		urlPath := input[:sepIndex]
		if len(urlPath) == 0 {
			continue
		}
		urlPath = util.CleanUrlPath(urlPath)

		fsPath := input[sepIndex+sepLen:]
		if len(fsPath) == 0 {
			continue
		}
		fsPath = path.Clean(fsPath)

		maps[urlPath] = fsPath
	}

	return maps
}

func normalizeUrlPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) > 0 {
			outputs = append(outputs, util.CleanUrlPath(input))
		}
	}

	return outputs
}
