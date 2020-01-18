package serverHandler

import (
	"bytes"
	"net/http"
)

const LOG_BUF_SIZE = 80

func (h *handler) logRequest(r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	payload := []byte(r.RemoteAddr + " " + r.Method + " " + r.RequestURI)

	h.logger.LogAccess(payload)
}

func (h *handler) logUpload(filename, fsPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buffer := bytes.NewBuffer(make([]byte, 0, LOG_BUF_SIZE))

	buffer.WriteString(r.RemoteAddr)
	buffer.WriteByte(' ')
	buffer.WriteString("save upload file: ")
	buffer.WriteString(filename)
	buffer.WriteString(" -> ")
	buffer.WriteString(fsPath)

	h.logger.LogAccess(buffer.Bytes())
}

func (h *handler) logArchive(filename, relPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buffer := bytes.NewBuffer(make([]byte, 0, LOG_BUF_SIZE))

	buffer.WriteString(r.RemoteAddr)
	buffer.WriteByte(' ')
	buffer.WriteString("archive file: \"")
	buffer.WriteString(filename)
	buffer.WriteString("\" <- ")
	buffer.WriteString(relPath)

	h.logger.LogAccess(buffer.Bytes())
}
