package serverHandler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h *aliasHandler) mkdirs(authUserName, fsPrefix string, files []string, aliasSubItems []os.FileInfo, r *http.Request) bool {
	var errs []error

	for _, inputFilename := range files {
		if len(inputFilename) == 0 {
			continue
		}

		var filename string
		var ok bool
		if filename, ok = getCleanDirFilePath(inputFilename); !ok {
			errs = append(errs, errors.New("mkdir: illegal directory path "+inputFilename))
			continue
		}

		filenamePart1 := filename
		if prefixEndIndex := strings.IndexByte(filenamePart1, '/'); prefixEndIndex > 0 {
			filenamePart1 = filenamePart1[0:prefixEndIndex]
		}
		if containsItem(aliasSubItems, filenamePart1) {
			errs = append(errs, errors.New("mkdir: ignore path shadowed by alias "+filename))
			continue
		}
		fsPath := filepath.Join(fsPrefix, filename)
		h.logMutate(authUserName, "mkdir", fsPath, r)
		err := os.MkdirAll(fsPath, 0755)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if h.logErrors(errs) {
		return false
	}

	return true
}
