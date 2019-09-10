package serverHandler

import (
	"../serverErrorHandler"
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"runtime"
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

	if size == 0 || f == nil {
		return nil
	}

	_, err = io.Copy(tw, f)
	if err != nil {
		return err
	}

	return nil
}

func (h *handler) tar(w http.ResponseWriter, r *http.Request, pageData *pageData) {
	tw := tar.NewWriter(w)
	defer func() {
		err := tw.Close()
		serverErrorHandler.LogError(err)
	}()

	filename := pageData.ItemName + ".tar"
	writeArchiveHeader(w, "application/octet-stream", filename)

	h.visitFs(
		h.root+pageData.handlerRequestPath,
		pageData.rawRequestPath,
		"",
		func(f *os.File, fInfo os.FileInfo, relPath string) {
			go h.logArchive(filename, relPath, r)
			err := writeTar(tw, f, fInfo, relPath)
			if serverErrorHandler.LogError(err) {
				runtime.Goexit()
			}
		},
	)
}

func (h *handler) tgz(w http.ResponseWriter, r *http.Request, pageData *pageData) {
	gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if serverErrorHandler.LogError(err) {
		return
	}
	defer func() {
		err := gzw.Close()
		serverErrorHandler.LogError(err)
	}()

	tw := tar.NewWriter(gzw)
	defer func() {
		err := tw.Close()
		serverErrorHandler.LogError(err)
	}()

	filename := pageData.ItemName + ".tar.gz"
	writeArchiveHeader(w, "application/octet-stream", filename)

	h.visitFs(
		h.root+pageData.handlerRequestPath,
		pageData.rawRequestPath,
		"",
		func(f *os.File, fInfo os.FileInfo, relPath string) {
			go h.logArchive(filename, relPath, r)
			err := writeTar(tw, f, fInfo, relPath)
			if serverErrorHandler.LogError(err) {
				runtime.Goexit()
			}
		},
	)
}
