package util

import (
	"io"
	"mime"
	"net/http"
	"path"
)

func GetContentType(filename string, rd io.Reader) (string, error) {
	ext := path.Ext(filename)
	ctype := mime.TypeByExtension(ext)
	if len(ctype) > 0 {
		return ctype, nil
	}

	var buf [512]byte
	n, err := rd.Read(buf[:])
	if err != nil {
		return ctype, err
	}

	ctype = http.DetectContentType(buf[:n])
	return ctype, nil
}
