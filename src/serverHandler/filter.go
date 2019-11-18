package serverHandler

import "os"

func (h *handler) FilterItems(items []os.FileInfo) []os.FileInfo {
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
		shouldShow := true
		if h.shows != nil {
			shouldShow = shouldShow && h.shows.MatchString(item.Name())
		}
		if h.showDirs != nil && item.IsDir() {
			shouldShow = shouldShow && h.showDirs.MatchString(item.Name())
		}
		if h.showFiles != nil && !item.IsDir() {
			shouldShow = shouldShow && h.showFiles.MatchString(item.Name())
		}

		shouldHide := false
		if h.hides != nil {
			shouldHide = shouldHide || h.hides.MatchString(item.Name())
		}
		if h.hideDirs != nil && item.IsDir() {
			shouldHide = shouldHide || h.hideDirs.MatchString(item.Name())
		}
		if h.hideFiles != nil && !item.IsDir() {
			shouldHide = shouldHide || h.hideFiles.MatchString(item.Name())
		}

		if shouldShow && !shouldHide {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
