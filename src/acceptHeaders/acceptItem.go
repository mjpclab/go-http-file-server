package acceptHeaders

import (
	"strconv"
	"strings"
)

const qualitySign = ";q="
const defaultQuality = 1000
const maxQualityDecimals = 3

type acceptItem struct {
	value   string
	quality int
}

func parseAcceptItem(input string) acceptItem {
	value := input
	if semiColonIndex := strings.IndexByte(value, ';'); semiColonIndex >= 0 {
		value = value[:semiColonIndex]
	}

	rest := input[len(value):]
	qSignIndex := strings.Index(rest, qualitySign)
	if qSignIndex < 0 {
		return acceptItem{value, defaultQuality}
	}

	rest = rest[qSignIndex+len(qualitySign):]
	if semiColonIndex := strings.IndexByte(rest, ';'); semiColonIndex >= 0 {
		rest = rest[:semiColonIndex]
	}
	qLen := len(rest)

	if qLen == 0 {
		return acceptItem{value, defaultQuality}
	}
	if qLen > 1 && rest[1] != '.' {
		return acceptItem{value, defaultQuality}
	}

	// "q=1" or q is an invalid value
	if rest[0] != '0' {
		return acceptItem{value, defaultQuality}
	}

	// "q=0."
	if qLen <= 2 {
		return acceptItem{value, 0}
	}

	rest = rest[2:]
	qDecimalLen := len(rest)
	if qDecimalLen > maxQualityDecimals {
		qDecimalLen = maxQualityDecimals
		rest = rest[:qDecimalLen]
	}

	quality, err := strconv.Atoi(rest)
	if err != nil {
		quality = defaultQuality
	} else {
		missingDigits := maxQualityDecimals - qDecimalLen
		for i := 0; i < missingDigits; i++ {
			quality *= 10
		}
	}
	return acceptItem{value, quality}
}
