package serverHandler

import (
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

	num := 0
	digits := 0
	for i, length := 0, len(filename); i < length; i++ {
		n := filename[i]
		if n >= '0' && n <= '9' {
			num *= 10
			num += int(n - '0')
			digits++
		} else {
			break
		}
	}

	filename = filename[digits:]
	for {
		num++
		newFilename := strconv.Itoa(num) + filename
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
