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

func (aliases aliases) byUrlPath(urlPath string) (aliasItem alias, ok bool) {
	for _, alias := range aliases {
		if alias.isMatch(urlPath) {
			return alias, true
		}
	}
	return alias{}, false
}

func (aliases aliases) Len() int {
	return len(aliases)
}

func (aliases aliases) Less(i, j int) bool {
	iLen := len(aliases[i].url)
	jLen := len(aliases[j].url)
	if iLen != jLen {
		// longer is prior
		return iLen > jLen
	}

	return i > j
}

func (aliases aliases) Swap(i, j int) {
	aliases[i], aliases[j] = aliases[j], aliases[i]
}
