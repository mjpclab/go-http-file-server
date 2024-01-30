package serverHandler

import (
	"sort"
)

type aliases []alias

func newAliases(entries [][2]string) aliases {
	aliases := make(aliases, 0, len(entries))
	for i := range entries {
		aliases = append(aliases, createAlias(entries[i][0], entries[i][1]))
	}

	sort.Sort(aliases)

	return aliases
}

func (list aliases) byUrlPath(urlPath string) (aliasItem alias, ok bool) {
	for _, alias := range list {
		if alias.isMatch(urlPath) {
			return alias, true
		}
	}
	return alias{}, false
}

func (list aliases) filterSuccessor(url string) aliases {
	var result aliases

	for _, a := range list {
		if a.isSuccessorOf(url) {
			result = append(result, a)
		}
	}

	return result
}

func (list aliases) Len() int {
	return len(list)
}

func (list aliases) Less(i, j int) bool {
	iLen := len(list[i].url)
	jLen := len(list[j].url)
	if iLen != jLen {
		// longer is prior
		return iLen > jLen
	}

	return i > j
}

func (list aliases) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
