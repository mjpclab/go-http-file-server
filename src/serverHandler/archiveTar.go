package serverHandler

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
)

func writeTar(tw *tar.Writer, f *os.File, fInfo os.FileInfo, archivePath string) error {
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

	header := &tar.Header{
		Name:       archivePath,
		Mode:       0664,
		Size:       size,
		ModTime:    fInfo.ModTime(),
		AccessTime: fInfo.ModTime(),
		ChangeTime: fInfo.ModTime(),
	}

	err := tw.WriteHeader(header)
	if err != nil {
		return err
	}

	if size == 0 || f == nil || fInfo.IsDir() {
		return nil
	}

	_, err = io.Copy(tw, f)
	if err != nil {
		return err
	}

	return nil
}

func (h *aliasHandler) tar(w http.ResponseWriter, r *http.Request, pageData *responseData) {
	if !pageData.CanArchive {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	selections, ok := h.normalizeArchiveSelections(r)
	if !ok {
		return
	}

	tw := tar.NewWriter(w)
	defer func() {
		err := tw.Close()
		h.logError(err)
	}()

	h.archive(
		w,
		r,
		pageData,
		selections,
		".tar",
		"application/octet-stream",
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			return writeTar(tw, f, fInfo, relPath)
		},
	)
}

func (h *aliasHandler) tgz(w http.ResponseWriter, r *http.Request, pageData *responseData) {
	if !pageData.CanArchive {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	selections, ok := h.normalizeArchiveSelections(r)
	if !ok {
		return
	}

	gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if h.logError(err) {
		return
	}
	defer func() {
		err := gzw.Close()
		h.logError(err)
	}()

	tw := tar.NewWriter(gzw)
	defer func() {
		err := tw.Close()
		h.logError(err)
	}()

	h.archive(
		w,
		r,
		pageData,
		selections,
		".tar.gz",
		"application/octet-stream",
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			return writeTar(tw, f, fInfo, relPath)
		},
	)
}
