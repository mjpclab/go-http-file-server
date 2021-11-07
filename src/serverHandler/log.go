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

func (h *handler) logMutate(username, action, detail string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buffer := bytes.NewBuffer(make([]byte, 0, LOG_BUF_SIZE))

	buffer.WriteString(r.RemoteAddr)
	if len(username) > 0 {
		buffer.WriteString(" (")
		buffer.WriteString(username)
		buffer.WriteByte(')')
	}
	buffer.WriteByte(' ')
	buffer.WriteString(action)
	buffer.WriteString(": ")
	buffer.WriteString(detail)

	h.logger.LogAccess(buffer.Bytes())
}

func (h *handler) logUpload(username, filename, fsPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buffer := bytes.NewBuffer(make([]byte, 0, LOG_BUF_SIZE))

	buffer.WriteString(r.RemoteAddr)
	if len(username) > 0 {
		buffer.WriteString(" (")
		buffer.WriteString(username)
		buffer.WriteByte(')')
	}
	buffer.WriteString(" upload: ")
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
	buffer.WriteString(" archive file: \"")
	buffer.WriteString(filename)
	buffer.WriteString("\" <- ")
	buffer.WriteString(relPath)

	h.logger.LogAccess(buffer.Bytes())
}
