package util

func Contains(collection []string, find string) bool {
	for i := range collection {
		if collection[i] == find {
			return true
		}
	}

	return false
}
