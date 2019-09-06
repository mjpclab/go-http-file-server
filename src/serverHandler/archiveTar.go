package serverHandler

import (
	"../serverError"
	"archive/tar"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
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
		serverError.LogError(err)
	}()

	filename := pageData.Item.Name()
	if filename == "." {
		filename = strings.Replace(r.Host, ":", "_", -1)
	}
	filename = url.PathEscape(filename + ".tar")

	header := w.Header()
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Disposition", "attachment; filename*=UTF-8''"+filename)
	header.Set("Cache-Control", "public, max-age=0")
	w.WriteHeader(http.StatusOK)

	h.visitFs(
		h.root+pageData.handlerRequestPath,
		pageData.rawRequestPath,
		"",
		func(f *os.File, fInfo os.FileInfo, relPath string) {
			go h.logArchive(filename, relPath, r)
			err := writeTar(tw, f, fInfo, relPath)
			if serverError.LogError(err) {
				runtime.Goexit()
			}
		},
	)
}
