package serverHandler

import (
	"net/http"
	"net/url"
)

const LOG_BUF_SIZE = 80

func (h *handler) logRequest(r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}
	uri, err := url.QueryUnescape(r.RequestURI)
	if err != nil {
		uri = r.RequestURI
	}

	payload := []byte(r.RemoteAddr + " " + r.Method + " " + uri)

	h.logger.LogAccess(payload)
}

func (h *handler) logMutate(username, action, detail string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := make([]byte, 0, LOG_BUF_SIZE)

	buf = append(buf, []byte(r.RemoteAddr)...) // 9-47 bytes, mainly 21 bytes
	if len(username) > 0 {
		buf = append(buf, []byte(" (")...) // 2 bytes
		buf = append(buf, []byte(username)...)
		buf = append(buf, ')') // 1 byte
	}
	buf = append(buf, ' ')               // 1 byte
	buf = append(buf, []byte(action)...) // 5-6 bytes
	buf = append(buf, []byte(": ")...)   // 2 bytes
	buf = append(buf, []byte(detail)...)

	h.logger.LogAccess(buf)
}

func (h *handler) logUpload(username, filename, fsPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := make([]byte, 0, LOG_BUF_SIZE)

	buf = append(buf, []byte(r.RemoteAddr)...) // 9-47 bytes, mainly 21 bytes
	if len(username) > 0 {
		buf = append(buf, []byte(" (")...) // 2 bytes
		buf = append(buf, []byte(username)...)
		buf = append(buf, ')') // 1 byte
	}
	buf = append(buf, []byte(" upload: ")...) // 9 bytes
	buf = append(buf, []byte(filename)...)
	buf = append(buf, []byte(" -> ")...) // 4 bytes
	buf = append(buf, []byte(fsPath)...)

	h.logger.LogAccess(buf)
}

func (h *handler) logArchive(filename, relPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := make([]byte, 0, LOG_BUF_SIZE)

	buf = append(buf, []byte(r.RemoteAddr)...)      // 9-47 bytes, mainly 21 bytes
	buf = append(buf, []byte(" archive file: ")...) // 15 bytes
	buf = append(buf, []byte(filename)...)
	buf = append(buf, []byte(" <- ")...) // 4 bytes
	buf = append(buf, []byte(relPath)...)

	h.logger.LogAccess(buf)
}
