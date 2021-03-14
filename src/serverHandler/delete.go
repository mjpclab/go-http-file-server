package serverHandler

import (
	"errors"
	"os"
)

func (h *handler) deleteItems(fsPrefix string, files []string, aliasSubItems []os.FileInfo) bool {
	errs := []error{}

	for _, inputFilename := range files {
		if len(inputFilename) == 0 {
			continue
		}

		var filename string
		var ok bool
		if filename, ok = getCleanFilePath(inputFilename); !ok {
			errs = append(errs, errors.New("delete: illegal item name "+inputFilename))
			continue
		}
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
		return false
	}

	return true
}
