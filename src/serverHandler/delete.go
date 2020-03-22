package serverHandler

import (
	"os"
)

func (h *handler) deleteItems(fsPrefix string, files []string) {
	errs := []error{}

	for _, filename := range files {
		err := os.Remove(fsPrefix + "/" + filename)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		go h.logger.LogErrors(errs...)
	}
}
