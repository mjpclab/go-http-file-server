package serverHandler

import (
	"os"
)

func (h *handler) mkdirs(fsPrefix string, files []string) {
	errs := []error{}

	for _, filename := range files {
		err := os.Mkdir(fsPrefix+"/"+filename, 0755)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		go h.logger.LogErrors(errs...)
	}
}
