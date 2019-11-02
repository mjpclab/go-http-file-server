package serverHandler

import (
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

func (h *handler) tar(w http.ResponseWriter, r *http.Request, pageData *responseData) {
	tw := tar.NewWriter(w)
	defer func() {
		err := tw.Close()
		h.errHandler.LogError(err)
	}()

	filename := pageData.ItemName + ".tar"
	writeArchiveHeader(w, "application/octet-stream", filename)

	if !needResponseBody(r.Method) {
		return
	}

	h.visitFs(
		h.root+pageData.handlerReqPath,
		pageData.rawReqPath,
		"",
		func(f *os.File, fInfo os.FileInfo, relPath string) {
			go h.logArchive(filename, relPath, r)
			err := writeTar(tw, f, fInfo, relPath)
			if h.errHandler.LogError(err) {
				runtime.Goexit()
			}
		},
	)
}

func (h *handler) tgz(w http.ResponseWriter, r *http.Request, pageData *responseData) {
	gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if h.errHandler.LogError(err) {
		return
	}
	defer func() {
		err := gzw.Close()
		h.errHandler.LogError(err)
	}()

	tw := tar.NewWriter(gzw)
	defer func() {
		err := tw.Close()
		h.errHandler.LogError(err)
	}()

	filename := pageData.ItemName + ".tar.gz"
	writeArchiveHeader(w, "application/octet-stream", filename)

	if !needResponseBody(r.Method) {
		return
	}

	h.visitFs(
		h.root+pageData.handlerReqPath,
		pageData.rawReqPath,
		"",
		func(f *os.File, fInfo os.FileInfo, relPath string) {
			go h.logArchive(filename, relPath, r)
			err := writeTar(tw, f, fInfo, relPath)
			if h.errHandler.LogError(err) {
				runtime.Goexit()
			}
		},
	)
}
