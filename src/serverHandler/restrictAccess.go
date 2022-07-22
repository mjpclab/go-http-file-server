package serverHandler

import (
	"../util"
	"net/http"
	"os"
	"strings"
)

func newRestrictAccesses(pathHostsMap map[string][]string) []pathStrings {
	restricts := make([]pathStrings, 0, len(pathHostsMap))

	for reqPath, hosts := range pathHostsMap {
		restricts = append(restricts, pathStrings{reqPath, hosts})
	}

	return restricts
}

func hasRestrictAccess(globalRestrictAccesses []string, restrictAccessUrls, restrictAccessDirs []pathStrings) bool {
	return globalRestrictAccesses != nil || len(restrictAccessUrls) > 0 || len(restrictAccessDirs) > 0
}

func (h *handler) isAllowAccess(r *http.Request, reqUrlPath, reqFsPath string, file *os.File, item os.FileInfo) bool {
	if !h.restrictAccess {
		return true
	}

	reqHeader := r.Header
	sourceHost := reqHeader.Get("Referer")
	if len(sourceHost) == 0 {
		sourceHost = reqHeader.Get("Origin")
	}

	if len(sourceHost) == 0 && !shouldServeAsContent(file, item) {
		return true
	}

	sourceHost = util.ExtractHostFromUrl(sourceHost)
	selfHost := strings.ToLower(r.Host)
	if sourceHost == selfHost {
		return true
	}

	if util.Contains(h.globalRestrictAccess, sourceHost) {
		return true
	}

	urlMatched := false
	for i := range h.restrictAccessUrls {
		if !util.HasUrlPrefixDir(reqUrlPath, h.restrictAccessUrls[i].path) {
			continue
		}
		urlMatched = true
		if util.Contains(h.restrictAccessUrls[i].strings, sourceHost) {
			return true
		}
	}

	dirMatched := false
	for i := range h.restrictAccessDirs {
		if !util.HasFsPrefixDir(reqFsPath, h.restrictAccessDirs[i].path) {
			continue
		}
		dirMatched = true
		if util.Contains(h.restrictAccessDirs[i].strings, sourceHost) {
			return true
		}
	}

	if h.globalRestrictAccess == nil && !urlMatched && !dirMatched {
		return true
	}

	return false
}

func restrictAccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("403 Forbidden"))
}
