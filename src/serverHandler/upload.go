package serverHandler

import (
	"../util"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const file = "file"
const dirFile = "dirfile"
const innerDirFile = "innerdirfile"

func getAvailableFilename(fsPrefix, filename string, mustAppendSuffix bool) string {
	if len(fsPrefix) == 0 {
		fsPrefix = "/"
	} else if fsPrefix[len(fsPrefix)-1] != '/' {
		fsPrefix = fsPrefix + "/"
	}

	if !mustAppendSuffix {
		if _, err := os.Lstat(fsPrefix + filename); os.IsNotExist(err) {
			return filename
		}
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

// RFC 7578, Section 4.2 requires that if a filename is provided, the
// directory path information must not be used.
// Since Go 1.17, Part.FileName() will strip directory information.
// However, the directory information is needed for uploading.
// Parse manually instead.
func getPartFilePath(part *multipart.Part) string {
	cd := part.Header.Get("Content-Disposition")
	_, params, _ := mime.ParseMediaType(cd)
	return params["filename"]
}

func (h *handler) saveUploadFiles(authUserName, fsPrefix string, createDir, overwriteExists bool, aliasSubItems []os.FileInfo, r *http.Request) bool {
	var errs []error

	reader, err := r.MultipartReader()
	if err != nil {
		errs = append(errs, err)
		return false
	}

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err != io.EOF {
				errs = append(errs, err)
			}
			break
		}

		inputPartFilePath := getPartFilePath(part)
		if len(inputPartFilePath) == 0 {
			continue
		}
		partFilePath, ok := getCleanDirFilePath(inputPartFilePath)
		if !ok {
			errs = append(errs, errors.New("upload: illegal file path "+inputPartFilePath))
			continue
		}

		filenameIndex := strings.LastIndexByte(partFilePath, '/')

		fsInfix := ""
		formname := part.FormName()
		if formname == dirFile {
			if filenameIndex > 0 {
				fsInfix = partFilePath[0:filenameIndex]
			}
		} else if formname == innerDirFile { // get file path, strip first level of dir
			if filenameIndex <= 0 {
				continue
			}
			filepath := partFilePath[0:filenameIndex]
			if prefixEndIndex := strings.IndexByte(filepath, '/'); prefixEndIndex > 0 {
				fsInfix = filepath[prefixEndIndex+1:]
			}
		} else if formname == file {
			// noop
		} else {
			errs = append(errs, errors.New("upload: unknown mode "+formname))
			continue
		}

		filePrefix := fsPrefix
		if len(fsInfix) > 0 {
			if !createDir {
				errs = append(errs, errors.New("upload: mkdir is not enabled for "+fsPrefix))
				continue
			}

			if len(aliasSubItems) > 0 {
				fsInfixPart1 := fsInfix
				fsInfixSlashIndex := strings.IndexByte(fsInfixPart1, '/')
				if fsInfixSlashIndex > 0 {
					fsInfixPart1 = fsInfixPart1[0:fsInfixSlashIndex]
				}
				if containsItem(aliasSubItems, fsInfixPart1) {
					errs = append(errs, errors.New("upload: ignore path shadowed by alias "+fsInfix))
					continue
				}
			}

			filePrefix += "/" + fsInfix
			err := os.MkdirAll(filePrefix, 0755)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		filename := partFilePath
		if filenameIndex >= 0 {
			filename = filename[filenameIndex+1:]
		}
		if len(filename) == 0 {
			continue
		}

		isFilenameAliased := len(fsInfix) == 0 && containsItem(aliasSubItems, filename)
		var fsFilename string
		if overwriteExists && !isFilenameAliased {
			tryPath := filePrefix + "/" + filename
			var info os.FileInfo
			info, err = os.Lstat(tryPath)
			if info != nil && !info.IsDir() {
				err = os.Remove(tryPath)
				if err != nil && !os.IsNotExist(err) {
					errs = append(errs, err)
				}
				// even remove failed, still try to write content to file by TRUNCATE mode
				fsFilename = filename
			}
		}
		if len(fsFilename) == 0 {
			fsFilename = getAvailableFilename(filePrefix, filename, isFilenameAliased)
		}
		if len(fsFilename) == 0 {
			err := errors.New("no available filename for " + filename)
			errs = append(errs, err)
			continue
		}

		fsPath := filepath.Join(filePrefix, fsFilename)
		h.logUpload(authUserName, filename, fsPath, r)
		file, err := os.OpenFile(fsPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		_, err = io.Copy(file, part)
		if err != nil {
			errs = append(errs, err)
		}

		err = file.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if h.logErrors(errs) {
		return false
	}

	return true
}
