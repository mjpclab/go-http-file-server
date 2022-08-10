package serverHandler

import (
	"../serverLog"
	"../shimgo"
	"../util"
	"net/http"
)

func (h *aliasHandler) logRequest(r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	var unescapedUri []byte
	unescapedLen := 0
	unescapedStr, err := shimgo.Net_Url_PathUnescape(r.RequestURI)
	if err == nil && unescapedStr != r.RequestURI {
		unescapedUri = util.EscapeControllingRune(unescapedStr)
		if len(unescapedUri) > 0 {
			unescapedLen = len(unescapedUri) + 5 // " <=> "
		}
	}

	uri := util.EscapeControllingRune(r.RequestURI)

	buf := serverLog.NewBuffer(2 + len(r.RemoteAddr) + len(r.Method) + unescapedLen + len(uri))

	buf = append(buf, []byte(r.RemoteAddr)...) // ~ 9-47 bytes, mainly 21 bytes
	buf = append(buf, ' ')                     // 1 byte
	buf = append(buf, []byte(r.Method)...)     // ~ 3-4 bytes
	buf = append(buf, ' ')                     // 1 byte
	if unescapedLen > 0 {
		buf = append(buf, unescapedUri...)
		buf = append(buf, ' ', '<', '=', '>', ' ') // 5 bytes
	}
	buf = append(buf, uri...)

	go h.logger.LogAccess(buf)
}

func (h *aliasHandler) logMutate(username, action, detail string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := serverLog.NewBuffer(6 + len(r.RemoteAddr) + len(username) + len(action) + len(detail))

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

	go h.logger.LogAccess(buf)
}

func (h *aliasHandler) logUpload(username, filename, fsPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := serverLog.NewBuffer(16 + len(r.RemoteAddr) + len(username) + len(filename) + len(fsPath))

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

	go h.logger.LogAccess(buf)
}

func (h *aliasHandler) logArchive(filename, relPath string, r *http.Request) {
	if !h.logger.CanLogAccess() {
		return
	}

	buf := serverLog.NewBuffer(19 + len(r.RemoteAddr) + len(filename) + len(relPath))

	buf = append(buf, []byte(r.RemoteAddr)...)      // ~ 9-47 bytes, mainly 21 bytes
	buf = append(buf, []byte(" archive file: ")...) // 15 bytes
	buf = append(buf, []byte(filename)...)
	buf = append(buf, []byte(" <- ")...) // 4 bytes
	buf = append(buf, []byte(relPath)...)

	go h.logger.LogAccess(buf)
}

func (h *aliasHandler) logErrors(errs []error) (hasError bool) {
	if len(errs) == 0 {
		return false
	}

	if h.logger.CanLogError() {
		go func(errs []error) {
			for i := range errs {
				errBytes := util.EscapeControllingRune(errs[i].Error())
				buf := serverLog.NewBuffer(len(errBytes))
				buf = append(buf, errBytes...)
				h.logger.LogError(buf)
			}
		}(errs)
	}

	return true
}

func (h *aliasHandler) logError(err error) (hasError bool) {
	if err == nil {
		return false
	}

	if h.logger.CanLogError() {
		go func(err error) {
			errBytes := util.EscapeControllingRune(err.Error())
			buf := serverLog.NewBuffer(len(errBytes))
			buf = append(buf, errBytes...)
			h.logger.LogError(buf)
		}(err)
	}

	return true
}
