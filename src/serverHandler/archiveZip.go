package serverHandler

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path"
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
	zipWriter := zip.NewWriter(w)
	defer func() {
		err := zipWriter.Close()
		h.errHandler.LogError(err)
	}()

	filename := pageData.ItemName + ".zip"
	writeArchiveHeader(w, "application/zip", filename)

	if !needResponseBody(r.Method) {
		return
	}

	h.visitFs(
		path.Clean(h.root+pageData.handlerReqPath),
		pageData.rawReqPath,
		"",
		func(f *os.File, fInfo os.FileInfo, relPath string) (err error) {
			go h.logArchive(filename, relPath, r)
			err = writeZip(zipWriter, f, fInfo, relPath)
			h.errHandler.LogError(err)
			return
		},
	)
}
