package acceptHeaders

import (
	"strconv"
	"strings"
)

const qualitySign = ";q="
const defaultQuality = 1000
const maxQualityDecimals = 3

type acceptItem struct {
	value     string
	quality   int
	wildcards int
}

func (item acceptItem) less(other acceptItem) bool {
	if item.quality != other.quality {
		return item.quality > other.quality
	}
	return item.wildcards < other.wildcards
}

func (item acceptItem) match(value string) bool {
	if item.value == value {
		return true
	}

	switch item.wildcards {
	case 0:
		return false
	case 1:
		slashIndex := strings.IndexByte(item.value, '/')
		itemPrefix := item.value[:slashIndex+1]
		return strings.HasPrefix(value, itemPrefix)
	case 2:
		return true
	}

	return false
}

func parseAcceptItem(input string) acceptItem {
	value := input
	if semiColonIndex := strings.IndexByte(value, ';'); semiColonIndex >= 0 {
		value = value[:semiColonIndex]
	}

	wildcards := 0
	if value == "*/*" {
		wildcards = 2
	} else if strings.HasSuffix(value, "/*") {
		wildcards = 1
	}

	rest := input[len(value):]
	qSignIndex := strings.Index(rest, qualitySign)
	if qSignIndex < 0 {
		return acceptItem{value, defaultQuality, wildcards}
	}

	rest = rest[qSignIndex+len(qualitySign):]
	if semiColonIndex := strings.IndexByte(rest, ';'); semiColonIndex >= 0 {
		rest = rest[:semiColonIndex]
	}
	qLen := len(rest)

	if qLen == 0 {
		return acceptItem{value, defaultQuality, wildcards}
	}
	if qLen > 1 && rest[1] != '.' {
		return acceptItem{value, defaultQuality, wildcards}
	}

	// "q=1" or q is an invalid value
	if rest[0] != '0' {
		return acceptItem{value, defaultQuality, wildcards}
	}

	// "q=0."
	if qLen <= 2 {
		return acceptItem{value, 0, wildcards}
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
	return acceptItem{value, quality, wildcards}
}
