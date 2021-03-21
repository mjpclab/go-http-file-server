package acceptHeaders

import (
	"strconv"
	"strings"
)

const qualitySign = ";q="
const defaultQuality = 1000
const maxQualityDigits = 3

type acceptItem struct {
	value   string
	quality int
}

func parseAcceptItem(input string) acceptItem {
	indexQSign := strings.Index(input, qualitySign)
	if indexQSign < 0 {
		return acceptItem{input, defaultQuality}
	}

	value := input[:indexQSign]
	strQuality := input[indexQSign+len(qualitySign):]
	if semiColonIndex := strings.IndexByte(strQuality, ';'); semiColonIndex >= 0 {
		strQuality = strQuality[:semiColonIndex]
	}
	strQualityLen := len(strQuality)

	if strQualityLen == 0 {
		return acceptItem{value, defaultQuality}
	}
	if strQualityLen > 1 && strQuality[1] != '.' {
		return acceptItem{value, defaultQuality}
	}

	// "q=1" or q is an invalid value
	if strQuality[0] != '0' {
		return acceptItem{value, defaultQuality}
	}

	// "q=0."
	if strQualityLen <= 2 {
		return acceptItem{value, 0}
	}

	strRest := strQuality[2:]
	strRestLen := len(strRest)
	if strRestLen > maxQualityDigits {
		strRestLen = maxQualityDigits
		strRest = strRest[:strRestLen]
	}

	quality, err := strconv.Atoi(strRest)
	if err != nil {
		quality = defaultQuality
	} else {
		missingDigits := maxQualityDigits - strRestLen
		for i := 0; i < missingDigits; i++ {
			quality *= 10
		}
	}
	return acceptItem{value, quality}
}
