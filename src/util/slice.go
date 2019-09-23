package util

func Contains(collection []string, find string) bool {
	for _, item := range collection {
		if item == find {
			return true
		}
	}

	return false
}
