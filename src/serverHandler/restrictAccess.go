package serverHandler

import (
	"net/http"
	"os"
	"strings"

	"mjpclab.dev/ghfs/src/util"
)

func newRestrictAccesses(pathHostsList [][]string) pathStringsList {
	restricts := make(pathStringsList, 0, len(pathHostsList))

	for _, pathHosts := range pathHostsList {
		if len(pathHosts) == 0 {
			continue
		}
		restricts = append(restricts, pathStrings{pathHosts[0], pathHosts[1:]})
	}

	return restricts
}

func (h *aliasHandler) isAllowAccess(r *http.Request, reqUrlPath, reqFsPath string, file *os.File, item os.FileInfo) (restrict, allow bool) {
	if h.globalRestrictAccess == nil && len(h.restrictAccessUrls) == 0 && len(h.restrictAccessDirs) == 0 {
		return false, true
	}

	reqHeader := r.Header
	sourceHost := reqHeader.Get("Referer")
	if len(sourceHost) == 0 {
		sourceHost = reqHeader.Get("Origin")
	}

	if len(sourceHost) == 0 && !shouldServeAsContent(file, item) {
		return true, true
	}

	sourceHost = util.ExtractHostFromUrl(sourceHost)
	selfHost := strings.ToLower(r.Host)
	if sourceHost == selfHost {
		return true, true
	}

	if util.Contains(h.globalRestrictAccess, sourceHost) {
		return true, true
	}

	urlMatched := false
	for i := range h.restrictAccessUrls {
		if !util.HasUrlPrefixDir(reqUrlPath, h.restrictAccessUrls[i].path) {
			continue
		}
		urlMatched = true
		if util.Contains(h.restrictAccessUrls[i].values, sourceHost) {
			return true, true
		}
	}

	dirMatched := false
	for i := range h.restrictAccessDirs {
		if !util.HasFsPrefixDir(reqFsPath, h.restrictAccessDirs[i].path) {
			continue
		}
		dirMatched = true
		if util.Contains(h.restrictAccessDirs[i].values, sourceHost) {
			return true, true
		}
	}

	if h.globalRestrictAccess == nil && !urlMatched && !dirMatched {
		return true, true
	}

	return true, false
}
