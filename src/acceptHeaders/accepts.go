package acceptHeaders

import (
	"sort"
	"strings"
)

type Accepts []acceptItem

func (accepts Accepts) Len() int {
	return len(accepts)
}

func (accepts Accepts) Swap(i, j int) {
	accepts[i], accepts[j] = accepts[j], accepts[i]
}

func (accepts Accepts) Less(i, j int) bool {
	return accepts[i].less(accepts[j])
}

// make sure `accepts` is sorted
func (accepts Accepts) GetPreferredValue(availables []string) (index int, value string, ok bool) {
	for _, accept := range accepts {
		for i, avail := range availables {
			if accept.match(avail) {
				return i, avail, true
			}
		}
	}

	index = -1
	return
}

func ParseAccepts(input string) Accepts {
	entries := strings.Split(input, ",")
	entryCount := len(entries)
	if entryCount == 0 {
		return nil
	}

	accepts := make(Accepts, 0, entryCount)
	for i := 0; i < entryCount; i++ {
		input := strings.TrimSpace(entries[i])
		if len(input) == 0 {
			continue
		}

		accept := parseAcceptItem(input)
		if accept.quality <= 0 {
			continue
		}

		accepts = append(accepts, accept)
	}
	sort.Sort(accepts)
	return accepts
}
