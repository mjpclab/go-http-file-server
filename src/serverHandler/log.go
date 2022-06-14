package serverHandler

import (
	"net/http"
)

func (h *handler) logRequest(r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := make([]byte, 0, 2+len(r.RemoteAddr)+len(r.Method)+len(r.RequestURI))

	buf = append(buf, []byte(r.RemoteAddr)...) // ~ 9-47 bytes, mainly 21 bytes
	buf = append(buf, ' ')                     // 1 byte
	buf = append(buf, []byte(r.Method)...)     // ~ 3-4 bytes
	buf = append(buf, ' ')                     // 1 byte
	buf = append(buf, []byte(r.RequestURI)...)

	h.logger.LogAccess(buf)
}

func (h *handler) logMutate(username, action, detail string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := make([]byte, 0, 6+len(r.RemoteAddr)+len(username)+len(action)+len(detail))

	buf = append(buf, []byte(r.RemoteAddr)...) // ~ 9-47 bytes, mainly 21 bytes
	if len(username) > 0 {
		buf = append(buf, ' ', '(') // 2 bytes
		buf = append(buf, []byte(username)...)
		buf = append(buf, ')') // 1 byte
	}
	buf = append(buf, ' ')               // 1 byte
	buf = append(buf, []byte(action)...) // ~ 5-6 bytes
	buf = append(buf, ':', ' ')          // 2 bytes
	buf = append(buf, []byte(detail)...)

	h.logger.LogAccess(buf)
}

func (h *handler) logUpload(username, filename, fsPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := make([]byte, 0, 16+len(r.RemoteAddr)+len(username)+len(filename)+len(fsPath))

	buf = append(buf, []byte(r.RemoteAddr)...) // ~ 9-47 bytes, mainly 21 bytes
	if len(username) > 0 {
		buf = append(buf, ' ', '(') // 2 bytes
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

	buf := make([]byte, 0, 19+len(r.RemoteAddr)+len(filename)+len(relPath))

	buf = append(buf, []byte(r.RemoteAddr)...)      // ~ 9-47 bytes, mainly 21 bytes
	buf = append(buf, []byte(" archive file: ")...) // 15 bytes
	buf = append(buf, []byte(filename)...)
	buf = append(buf, []byte(" <- ")...) // 4 bytes
	buf = append(buf, []byte(relPath)...)

	h.logger.LogAccess(buf)
}
