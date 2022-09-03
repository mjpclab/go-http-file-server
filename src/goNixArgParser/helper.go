package goNixArgParser

import (
	"os"
	"regexp"
)

// ==================================================

func removeEmptyInplace(inputs []string) (outputs []string) {
	// remove empty item
	nextIndex := 0
	for i := range inputs {
		if len(inputs[i]) == 0 {
			continue
		}
		if nextIndex != i {
			inputs[nextIndex] = inputs[i]
		}
		nextIndex++
	}

	return inputs[:nextIndex]
}

// ==================================================

/*
- non-quoted: [^'"\s]*
- repeat of quoted and non-quoted: (?:(?:'[^']*'|"[^"]*")[^'"\s]*)*
- rest unpaired quote: \S*
*/
var reLexical = regexp.MustCompile(`[^'"\s]*(?:(?:'[^']*'|"[^"]*")[^'"\s]*)*\S*`)

func splitLexicals(input string) (output []string) {
	lexicals := reLexical.FindAllString(input, -1)
	lexicals = removeEmptyInplace(lexicals)

	return lexicals
}

// ==================================================

var reQuoted = regexp.MustCompile(`'[^']*'|"[^"]*"`)

func removeQuotes(input string) (output string) {
	output = reQuoted.ReplaceAllStringFunc(input, func(matched string) string {
		return matched[1 : len(matched)-1]
	})
	return output
}

func removeAllQuotesInPlace(inputs []string) []string {
	for i := range inputs {
		inputs[i] = removeQuotes(inputs[i])
	}
	inputs = removeEmptyInplace(inputs)
	return inputs
}

// ==================================================

func SplitToArgs(input string) (args []string) {
	lexicals := splitLexicals(input)
	args = removeAllQuotesInPlace(lexicals)
	return args
}

func LoadConfigArgs(filename string) (args []string, err error) {
	input, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	strConfig := string(input)
	return SplitToArgs(strConfig), nil
}
