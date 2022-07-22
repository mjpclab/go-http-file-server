package serverHandler

import (
	"sort"
)

type aliases []alias

func newAliases(entries map[string]string) aliases {
	aliases := make(aliases, 0, len(entries))
	for urlPath, fsPath := range entries {
		aliases = append(aliases, createAlias(urlPath, fsPath))
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

	return i < j
}

func (aliases aliases) Swap(i, j int) {
	aliases[i], aliases[j] = aliases[j], aliases[i]
}
