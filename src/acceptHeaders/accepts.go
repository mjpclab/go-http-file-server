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
	return accepts[i].quality > accepts[j].quality
}

// make sure `accepts` is sorted
func (accepts Accepts) GetPreferredValue(availables []string) (index int, value string, ok bool) {
	for _, accept := range accepts {
		for i, avail := range availables {
			if accept.value == avail {
				return i, avail, true
			}
		}
	}
	return
}

func ParseAccepts(input string) Accepts {
	entries := strings.Split(input, ",")
	entryCount := len(entries)
	if entryCount == 0 {
		return nil
	}

	accepts := make(Accepts, entryCount)
	for i := 0; i < entryCount; i++ {
		accepts[i] = parseAcceptItem(strings.TrimSpace(entries[i]))
	}
	sort.Sort(accepts)
	return accepts
}
