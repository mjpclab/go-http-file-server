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

func stringToSlice(input string) []string {
	if len(input) == 0 {
		return nil
	}

	return []string{input}
}

func toBool(input string) (bool, bool) {
	if len(input) == 0 {
		return false, false
	}
	v, err := strconv.ParseBool(input)
	return v, err == nil
}

func toBools(inputs []string) (outputs []bool, ok bool) {
	l := len(inputs)
	values := make([]bool, l)
	for i := 0; i < l; i++ {
		values[i], ok = toBool(inputs[i])
		if !ok {
			return
		}
	}

	return values, true
}

func toInt(input string) (int, bool) {
	v, err := strconv.Atoi(input)
	return v, err == nil
}

func toInts(inputs []string) (outputs []int, ok bool) {
	l := len(inputs)
	values := make([]int, l)
	for i := 0; i < l; i++ {
		values[i], ok = toInt(inputs[i])
		if !ok {
			return
		}
	}

	return values, true
}

func toUint(input string) (uint, bool) {
	v, err := strconv.ParseUint(input, 10, 0)
	return uint(v), err == nil
}

func toUints(inputs []string) (outputs []uint, ok bool) {
	l := len(inputs)
	values := make([]uint, l)
	for i := 0; i < l; i++ {
		values[i], ok = toUint(inputs[i])
		if !ok {
			return
		}
	}

	return values, true
}

func toInt32(input string) (int32, bool) {
	v, err := strconv.ParseInt(input, 10, 32)
	return int32(v), err == nil
}

func toInt32s(inputs []string) (outputs []int32, ok bool) {
	l := len(inputs)
	values := make([]int32, l)
	for i := 0; i < l; i++ {
		values[i], ok = toInt32(inputs[i])
		if !ok {
			return
		}
	}

	return values, true
}

func toUint32(input string) (uint32, bool) {
	v, err := strconv.ParseUint(input, 10, 32)
	return uint32(v), err == nil
}

func toUint32s(inputs []string) (outputs []uint32, ok bool) {
	l := len(inputs)
	values := make([]uint32, l)
	for i := 0; i < l; i++ {
		values[i], ok = toUint32(inputs[i])
		if !ok {
			return
		}
	}

	return values, true
}

func toInt64(input string) (int64, bool) {
	v, err := strconv.ParseInt(input, 10, 64)
	return v, err == nil
}

func toInt64s(inputs []string) (outputs []int64, ok bool) {
	l := len(inputs)
	values := make([]int64, l)
	for i := 0; i < l; i++ {
		values[i], ok = toInt64(inputs[i])
		if !ok {
			return
		}
	}

	return values, true
}

func toUint64(input string) (uint64, bool) {
	v, err := strconv.ParseUint(input, 10, 64)
	return v, err == nil
}

func toUint64s(inputs []string) (outputs []uint64, ok bool) {
	l := len(inputs)
	values := make([]uint64, l)
	for i := 0; i < l; i++ {
		values[i], ok = toUint64(inputs[i])
		if !ok {
			return
		}
	}

	return values, true
}

func toFloat64(input string) (float64, bool) {
	v, err := strconv.ParseFloat(input, 64)
	return v, err == nil
}

func toFloat64s(inputs []string) (outputs []float64, ok bool) {
	l := len(inputs)
	values := make([]float64, l)
	for i := 0; i < l; i++ {
		values[i], ok = toFloat64(inputs[i])
		if !ok {
			return
		}
	}

	return values, true
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
