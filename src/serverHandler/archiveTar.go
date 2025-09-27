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

	var typeFlag byte
	var mode int64
	var size int64
	if fInfo.IsDir() {
		archivePath += "/"
		typeFlag = tar.TypeDir
		mode = 0755
	} else {
		typeFlag = tar.TypeReg
		mode = 0644
		size = fInfo.Size()
	}

	header := &tar.Header{
		Name:       archivePath,
		Typeflag:   typeFlag,
		Mode:       mode,
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

func (h *aliasHandler) tar(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) bool {
	if !data.CanArchive {
		data.Status = http.StatusBadRequest
		return false
	}

	selections, ok := h.normalizeArchiveSelections(r)
	if !ok {
		data.Status = http.StatusBadRequest
		return false
	}

	tw := tar.NewWriter(w)
	defer func() {
		err := tw.Close()
		h.logError(err)
	}()

	h.archiveFiles(
		w,
		r,
		session,
		data,
		selections,
		"application/x-tar",
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			return writeTar(tw, f, fInfo, relPath)
		},
	)
	return true
}

func (h *aliasHandler) tgz(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) bool {
	if !data.CanArchive {
		data.Status = http.StatusBadRequest
		return false
	}

	selections, ok := h.normalizeArchiveSelections(r)
	if !ok {
		data.Status = http.StatusBadRequest
		return false
	}

	gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if h.logError(err) {
		data.Status = http.StatusInternalServerError
		return false
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

	h.archiveFiles(
		w,
		r,
		session,
		data,
		selections,
		"application/gzip",
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			return writeTar(tw, f, fInfo, relPath)
		},
	)
	return true
}
