package serverHandler

import "os"

func (h *aliasHandler) FilterItems(items []os.FileInfo) []os.FileInfo {
	if h.shows == nil &&
		h.showDirs == nil &&
		h.showFiles == nil &&
		h.hides == nil &&
		h.hideDirs == nil &&
		h.hideFiles == nil {
		return items
	}

	filtered := make([]os.FileInfo, 0, len(items))

	for _, item := range items {
		name := item.Name()

		if h.hides != nil && h.hides.MatchString(name) {
			continue
		}

		if h.hideDirs != nil && item.IsDir() && h.hideDirs.MatchString(name) {
			continue
		}

		if h.hideFiles != nil && !item.IsDir() && h.hideFiles.MatchString(name) {
			continue
		}

		if h.shows != nil && !h.shows.MatchString(name) {
			continue
		}

		if h.showDirs != nil && item.IsDir() && !h.showDirs.MatchString(name) {
			continue
		}

		if h.showFiles != nil && !item.IsDir() && !h.showFiles.MatchString(name) {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}
