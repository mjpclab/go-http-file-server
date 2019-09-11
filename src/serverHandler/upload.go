package serverHandler

import (
	"../util"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

func getAvailableFilename(fsPrefix, filename string) string {
	if len(fsPrefix) == 0 {
		fsPrefix = "/"
	} else if fsPrefix[len(fsPrefix)-1] != '/' {
		fsPrefix = fsPrefix + "/"
	}

	if _, err := os.Lstat(fsPrefix + filename); os.IsNotExist(err) {
		return filename
	}

	filenamePrefix, filenameSuffix := util.SplitFilename(filename)

	for i := 1; ; i++ {
		newFilename := filenamePrefix + "-" + strconv.Itoa(i) + filenameSuffix
		if _, err := os.Lstat(fsPrefix + newFilename); os.IsNotExist(err) {
			return newFilename
		}
	}

	return ""
}

func (h *handler) saveUploadFiles(requestPath string, r *http.Request) (errs []error) {
	errs = []error{}

	reader, err := r.MultipartReader()
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := 1; ; i++ {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}

		filename := part.FileName()
		if len(filename) == 0 {
			continue
		}
		fsFilename := getAvailableFilename(h.root+requestPath, filename)
		if len(fsFilename) == 0 {
			err := errors.New("no available filename for " + filename)
			errs = append(errs, err)
			continue
		}

		fsPath := path.Clean(h.root + requestPath + "/" + fsFilename)
		go h.logUpload(filename, fsPath, r)
		file, err := os.Create(fsPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		_, err = io.Copy(file, part)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = file.Close()
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return
}
