package acceptHeaders

import (
	"mjpclab.dev/ghfs/src/util"
	"strconv"
	"strings"
)

const qualitySign = "q="
const defaultQuality = 1000
const maxQualityDecimals = 3

type acceptItem struct {
	value   string
	quality int
}

func (item acceptItem) less(other acceptItem) bool {
	if item.quality != other.quality {
		return item.quality > other.quality
	}

	if item.value != other.value {
		if other.value == "*/*" {
			return true
		} else if strings.HasSuffix(other.value, "/*") && !strings.HasPrefix(item.value, "/*") {
			return true
		}
	}

	return false
}

func (item acceptItem) match(value string) bool {
	if item.value == value {
		return true
	}
	if item.value == "*/*" {
		return true
	}
	if !strings.HasSuffix(item.value, "/*") {
		return false
	}

	itemPrefix := item.value[:len(item.value)-1]
	return strings.HasPrefix(value, itemPrefix)
}

func parseAcceptItem(input string) acceptItem {
	entries := strings.Split(input, ";")
	if len(entries) == 0 {
		return acceptItem{}
	}
	util.InPlaceTrim(entries)

	value := entries[0]
	entries = entries[1:]

	quality := defaultQuality
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if strings.HasPrefix(entry, qualitySign) {
			quality = parseQuality(entry[len(qualitySign):])
		}
	}

	return acceptItem{value, quality}
}

func parseQuality(input string) (quality int) {
	qLen := len(input)

	if qLen == 0 {
		return defaultQuality
	}
	if qLen > 1 && input[1] != '.' {
		return defaultQuality
	}

	// q is "1" or q is an invalid value
	if input[0] != '0' {
		return defaultQuality
	}

	// "0."
	if qLen <= 2 {
		return 0
	}

	input = input[2:]
	qDecimalLen := len(input)
	if qDecimalLen > maxQualityDecimals {
		qDecimalLen = maxQualityDecimals
		input = input[:qDecimalLen]
	}

	quality, err := strconv.Atoi(input)
	if err != nil {
		quality = defaultQuality
	} else {
		missingDigits := maxQualityDecimals - qDecimalLen
		for i := 0; i < missingDigits; i++ {
			quality *= 10
		}
	}
	return quality
}
