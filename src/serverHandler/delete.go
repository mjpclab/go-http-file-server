package serverHandler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
)

func (h *aliasHandler) deleteItems(authUserName, fsPrefix string, files []string, aliasSubItems []os.FileInfo, r *http.Request) bool {
	var errs []error

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
		fsPath := filepath.Join(fsPrefix, filename)
		h.logMutate(authUserName, "delete", fsPath, r)
		err := os.RemoveAll(fsPath)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if h.logErrors(errs) {
		return false
	}

	return true
}
