package serverHandler

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
)

func writeZip(zw *zip.Writer, f *os.File, fInfo os.FileInfo, archivePath string) error {
	if archivePath[0] == '/' {
		archivePath = archivePath[1:]
	}
	if fInfo.IsDir() {
		archivePath += "/"
	}

	var size int64
	if !fInfo.IsDir() {
		size = fInfo.Size()
	}

	w, err := zw.Create(archivePath)
	if err != nil {
		return err
	}

	if size == 0 || f == nil || fInfo.IsDir() {
		return nil
	}

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	return nil
}

func (h *aliasHandler) zip(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) bool {
	if !data.CanArchive {
		data.Status = http.StatusBadRequest
		return false
	}

	selections, ok := h.normalizeArchiveSelections(r)
	if !ok {
		data.Status = http.StatusBadRequest
		return false
	}

	zipWriter := zip.NewWriter(w)
	defer func() {
		err := zipWriter.Close()
		h.logError(err)
	}()

	h.archiveFiles(
		w,
		r,
		session,
		data,
		selections,
		"application/zip",
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			return writeZip(zipWriter, f, fInfo, relPath)
		},
	)
	return true
}
