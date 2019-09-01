package serverHandler

import (
	"bytes"
	"net/http"
)

func (h *handler) logRequest(r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buffer := &bytes.Buffer{}

	buffer.WriteString(r.RemoteAddr)
	buffer.WriteByte(' ')
	buffer.WriteString(r.Method)
	buffer.WriteByte(' ')
	buffer.WriteString(r.RequestURI)

	h.logger.LogAccess(buffer.Bytes())
}

func (h *handler) logUpload(filename, fsPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buffer := &bytes.Buffer{}

	buffer.WriteString(r.RemoteAddr)
	buffer.WriteByte(' ')
	buffer.WriteString("save upload file: ")
	buffer.WriteString(filename)
	buffer.WriteString(" -> ")
	buffer.WriteString(fsPath)

	h.logger.LogAccess(buffer.Bytes())
}
