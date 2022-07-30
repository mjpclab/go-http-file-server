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

func (h *aliasHandler) getCanUpload(info os.FileInfo, rawReqPath, reqFsPath string) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.globalUpload {
		return true
	}

	return hasUrlOrDirPrefix(h.uploadUrls, rawReqPath, h.uploadDirs, reqFsPath)
}

func (h *aliasHandler) getCanMkdir(info os.FileInfo, rawReqPath, reqFsPath string) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.globalMkdir {
		return true
	}

	return hasUrlOrDirPrefix(h.mkdirUrls, rawReqPath, h.mkdirDirs, reqFsPath)
}

func (h *aliasHandler) getCanDelete(info os.FileInfo, rawReqPath, reqFsPath string) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.globalDelete {
		return true
	}

	return hasUrlOrDirPrefix(h.deleteUrls, rawReqPath, h.deleteDirs, reqFsPath)
}

func (h *aliasHandler) getCanArchive(subInfos []os.FileInfo, rawReqPath, reqFsPath string) bool {
	if len(subInfos) == 0 {
		return false
	}

	if h.globalArchive {
		return true
	}

	return hasUrlOrDirPrefix(h.archiveUrls, rawReqPath, h.archiveDirs, reqFsPath)
}

func (h *aliasHandler) getCanCors(rawReqPath, reqFsPath string) bool {
	if h.globalCors {
		return true
	}

	return hasUrlOrDirPrefix(h.corsUrls, rawReqPath, h.corsDirs, reqFsPath)
}

func (h *aliasHandler) getNeedAuth(rawReqPath, reqFsPath string) bool {
	if h.globalAuth {
		return true
	}

	return hasUrlOrDirPrefix(h.authUrls, rawReqPath, h.authDirs, reqFsPath)
}
