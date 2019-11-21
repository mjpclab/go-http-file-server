package serverHandler

import (
	"../util"
	"os"
)

func hasUrlOrDirPrefix(urls []string, reqUrl string, dirs []string, reqDir string) bool {
	for _, url := range urls {
		if util.HasUrlPrefixDir(reqUrl, url) {
			return true
		}
	}

	for _, dir := range dirs {
		if util.HasFsPrefixDir(reqDir, dir) {
			return true
		}
	}

	return false
}

func (h *handler) getCanUpload(item os.FileInfo, rawReqPath, reqFsPath string) bool {
	if item == nil || !item.IsDir() {
		return false
	}

	if h.globalUpload {
		return true
	}

	return hasUrlOrDirPrefix(h.uploadUrls, rawReqPath, h.uploadDirs, reqFsPath)
}

func (h *handler) getCanArchive(subItems []os.FileInfo, rawReqPath, reqFsPath string) bool {
	if len(subItems) == 0 {
		return false
	}

	if h.globalArchive {
		return true
	}

	return hasUrlOrDirPrefix(h.archiveUrls, rawReqPath, h.archiveDirs, reqFsPath)
}

func (h *handler) getCanCors(rawReqPath, reqFsPath string) bool {
	if h.globalCors {
		return true
	}

	return hasUrlOrDirPrefix(h.corsUrls, rawReqPath, h.corsDirs, reqFsPath)
}

func (h *handler) getNeedAuth(rawReqPath, reqFsPath string) bool {
	if h.globalAuth {
		return true
	}

	return hasUrlOrDirPrefix(h.authUrls, rawReqPath, h.authDirs, reqFsPath)
}
