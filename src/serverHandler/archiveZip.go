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

func (h *handler) zip(w http.ResponseWriter, r *http.Request, pageData *responseData) {
	if !pageData.CanArchive {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	selections, ok := h.normalizeArchiveSelections(r)
	if !ok {
		return
	}

	zipWriter := zip.NewWriter(w)
	defer func() {
		err := zipWriter.Close()
		h.logError(err)
	}()

	h.archive(
		w,
		r,
		pageData,
		selections,
		".zip",
		"application/zip",
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			return writeZip(zipWriter, f, fInfo, relPath)
		},
	)
}
