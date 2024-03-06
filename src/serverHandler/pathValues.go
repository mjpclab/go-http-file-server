package serverHandler

// prefixFilter

type prefixFilter func(whole, prefix string) bool

// pathValues

type pathValues[T any] struct {
	path   string
	values []T
}

type pathValuesList[T any] []pathValues[T]

func (list pathValuesList[T]) mergePrefixMatched(mergeWith []T, matchPrefix prefixFilter, refPath string) []T {
	var result []T
	if mergeWith != nil {
		result = make([]T, len(mergeWith))
		copy(result, mergeWith)
	}

	for i := range list {
		if matchPrefix(refPath, list[i].path) {
			if result == nil {
				result = []T{}
			}
			result = append(result, list[i].values...)
		}
	}

	if mergeWith != nil && len(mergeWith) == len(result) {
		return mergeWith
	} else {
		return result
	}
}

func (list pathValuesList[T]) filterSuccessor(includeSelf bool, matchPrefix prefixFilter, refPath string) pathValuesList[T] {
	var result pathValuesList[T]

	for i := range list {
		if !includeSelf && len(list[i].path) == len(refPath) {
			continue
		}
		if matchPrefix(list[i].path, refPath) {
			result = append(result, list[i])
		}
	}

	if len(list) == len(result) {
		return list
	} else {
		return result
	}
}

// pathInts

type pathInts = pathValues[int]
type pathIntsList = pathValuesList[int]

// pathStrings

type pathStrings = pathValues[string]
type pathStringsList = pathValuesList[string]

// pathHeaders

type pathHeaders = pathValues[[2]string]
type pathHeadersList = pathValuesList[[2]string]

// []string

func prefixMatched(list []string, matchPrefix prefixFilter, refPath string) bool {
	for i := range list {
		if matchPrefix(refPath, list[i]) {
			return true
		}
	}

	return false
}

func filterSuccessor(list []string, matchPrefix prefixFilter, refPath string) []string {
	var result []string

	for _, v := range list {
		if len(v) > len(refPath) && matchPrefix(v, refPath) {
			result = append(result, v)
		}
	}

	if len(list) == len(result) {
		return list
	} else {
		return result
	}
}
