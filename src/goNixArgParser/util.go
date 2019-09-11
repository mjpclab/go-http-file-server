package goNixArgParser

import "strconv"

func getValue(source map[string][]string, key string) (value string, found bool) {
	var values []string
	values, found = source[key]

	if found && len(values) > 0 {
		value = values[0]
	}

	return
}

func getValues(source map[string][]string, key string) (values []string, found bool) {
	values, found = source[key]
	if found {
		values = copys(values)
		return values, true
	}
	return
}

func copys(input []string) []string {
	if input == nil {
		return nil
	}

	output := make([]string, len(input))
	copy(output, input)
	return output
}

func toBool(input string) (bool, error) {
	if len(input) == 0 {
		return false, nil
	}
	return strconv.ParseBool(input)
}

func toBools(input []string) ([]bool, error) {
	inputLen := len(input)

	output := make([]bool, inputLen)
	for i, l := 0, inputLen; i < l; i++ {
		v, err := toBool(input[i])
		if err != nil {
			return nil, err
		}
		output[i] = v
	}

	return output, nil
}

func toInt(input string) (int, error) {
	return strconv.Atoi(input)
}

func toInts(input []string) ([]int, error) {
	inputLen := len(input)

	output := make([]int, inputLen)
	for i, l := 0, inputLen; i < l; i++ {
		v, err := toInt(input[i])
		if err != nil {
			return nil, err
		}
		output[i] = v
	}

	return output, nil
}

func toInt64(input string) (int64, error) {
	return strconv.ParseInt(input, 10, 64)
}

func toInt64s(input []string) ([]int64, error) {
	inputLen := len(input)

	output := make([]int64, inputLen)
	for i, l := 0, inputLen; i < l; i++ {
		v, err := toInt64(input[i])
		if err != nil {
			return nil, err
		}
		output[i] = v
	}

	return output, nil
}

func toUint64(input string) (uint64, error) {
	return strconv.ParseUint(input, 10, 64)
}

func toUint64s(input []string) ([]uint64, error) {
	inputLen := len(input)

	output := make([]uint64, inputLen)
	for i, l := 0, inputLen; i < l; i++ {
		v, err := toUint64(input[i])
		if err != nil {
			return nil, err
		}
		output[i] = v
	}

	return output, nil
}

func toFloat64(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}

func toFloat64s(input []string) ([]float64, error) {
	inputLen := len(input)

	output := make([]float64, inputLen)
	for i, l := 0, inputLen; i < l; i++ {
		v, err := toFloat64(input[i])
		if err != nil {
			return nil, err
		}
		output[i] = v
	}

	return output, nil
}

func contains(collection []string, find string) bool {
	for _, item := range collection {
		if item == find {
			return true
		}
	}
	return false
}

func appendUnique(origins []string, items ...string) []string {
	for _, item := range items {
		if !contains(origins, item) {
			origins = append(origins, item)
		}
	}

	return origins
}
