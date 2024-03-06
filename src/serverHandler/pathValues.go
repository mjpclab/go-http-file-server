package serverHandler

// prefixFilter

type prefixFilter func(whole, prefix string) bool

// pathInts

type pathInts struct {
	path   string
	values []int
}

type pathIntsList []pathInts

func (list pathIntsList) mergePrefixMatched(mergeWith []int, matchPrefix prefixFilter, refPath string) []int {
	var result []int
	if mergeWith != nil {
		result = make([]int, len(mergeWith))
		copy(result, mergeWith)
	}

	for i := range list {
		if matchPrefix(refPath, list[i].path) {
			if result == nil {
				result = []int{}
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

func (list pathIntsList) filterSuccessor(includeSelf bool, matchPrefix prefixFilter, refPath string) pathIntsList {
	var result pathIntsList

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

// pathStrings

type pathStrings struct {
	path   string
	values []string
}

type pathStringsList []pathStrings

func (list pathStringsList) mergePrefixMatched(mergeWith []string, matchPrefix prefixFilter, refPath string) []string {
	var result []string
	if mergeWith != nil {
		result = make([]string, len(mergeWith))
		copy(result, mergeWith)
	}

	for i := range list {
		if matchPrefix(refPath, list[i].path) {
			if result == nil {
				result = []string{}
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

func (list pathStringsList) filterSuccessor(includeSelf bool, matchPrefix prefixFilter, refPath string) pathStringsList {
	var result pathStringsList

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

// pathHeaders

type pathHeaders struct {
	path   string
	values [][2]string
}

type pathHeadersList []pathHeaders

func (list pathHeadersList) mergePrefixMatched(mergeWith [][2]string, matchPrefix prefixFilter, refPath string) [][2]string {
	var result [][2]string
	if mergeWith != nil {
		result = make([][2]string, len(mergeWith))
		copy(result, mergeWith)
	}

	for i := range list {
		if matchPrefix(refPath, list[i].path) {
			if result == nil {
				result = [][2]string{}
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

func (list pathHeadersList) filterSuccessor(includeSelf bool, matchPrefix prefixFilter, refPath string) pathHeadersList {
	var result pathHeadersList

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
