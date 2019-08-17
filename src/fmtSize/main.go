package fmtSize

import (
	"strconv"
	"strings"
)

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
)

func fmtUnit(unitName string, unitValue int64, srcValue int64) string {
	prefix := int(srcValue / unitValue)
	suffix := int(srcValue % unitValue * 100 / unitValue)

	if suffix >= 55 {
		prefix++
	}

	b := strings.Builder{}
	b.WriteString(strconv.Itoa(prefix))
	b.WriteString(unitName)

	return b.String()
}

func FmtSize(size int64) string {
	switch {
	case size > PB:
		return fmtUnit("P", PB, size)
	case size > TB:
		return fmtUnit("T", TB, size)
	case size > GB:
		return fmtUnit("G", GB, size)
	case size > MB:
		return fmtUnit("M", MB, size)
	case size > KB:
		return fmtUnit("K", KB, size)
	default:
		return fmtUnit("", B, size)
	}
}
