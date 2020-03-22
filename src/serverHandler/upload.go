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

func (h *handler) saveUploadFiles(fsPrefix string, r *http.Request) {
	errs := []error{}

	reader, err := r.MultipartReader()
	if err != nil {
		errs = append(errs, err)
		return
	}

	for {
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
		fsFilename := getAvailableFilename(fsPrefix, filename)
		if len(fsFilename) == 0 {
			err := errors.New("no available filename for " + filename)
			errs = append(errs, err)
			continue
		}

		fsPath := path.Clean(fsPrefix + "/" + fsFilename)
		go h.logUpload(filename, fsPath, r)
		file, err := os.OpenFile(fsPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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

	if len(errs) > 0 {
		go h.logger.LogErrors(errs...)
	}
}
