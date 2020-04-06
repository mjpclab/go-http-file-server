package serverHandler

import (
	"os"
)

func (h *handler) deleteItems(fsPrefix string, files []string, aliasSubItems []os.FileInfo) {
	errs := []error{}

	for _, filename := range files {
		if containsItem(aliasSubItems, filename) {
			continue
		}
		err := os.RemoveAll(fsPrefix + "/" + filename)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		go h.logger.LogErrors(errs...)
	}
}
