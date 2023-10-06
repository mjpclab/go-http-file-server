package util

import (
	"os"
	"strings"
)

func GetBoolEnv(key string) bool {
	return isTruthyEnvValue(os.Getenv(key))
}

func isTruthyEnvValue(value string) bool {
	valueLen := len(value)

	if valueLen == 0 {
		return false
	}

	if strings.EqualFold(value, "false") {
		return false
	}

	for i := 0; i < valueLen; i++ {
		if value[i] != '0' &&
			value[i] != ' ' &&
			value[i] != '\t' &&
			value[i] != '\f' &&
			value[i] != '\r' &&
			value[i] != '\n' {
			return true
		}
	}

	return false
}
