package serverHandler

// prefixFilter

type prefixFilter func(whole, prefix string) bool

// pathStrings

type pathStrings struct {
	path    string
	strings []string
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
			result = append(result, list[i].strings...)
		}
	}

	if mergeWith != nil && len(mergeWith) == len(result) {
		return mergeWith
	} else {
		return result
	}
}

func (list pathStringsList) filterSuccessor(matchPrefix prefixFilter, refPath string) pathStringsList {
	var result pathStringsList

	for i := range list {
		if len(list[i].path) > len(refPath) && matchPrefix(list[i].path, refPath) {
			result = append(result, list[i])
		}
	}

	if len(list) == len(result) {
		return list
	} else {
		return result
	}
}
