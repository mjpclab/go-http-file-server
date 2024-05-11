package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
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

func hasUrlOrDirPrefixUsers(urlsUsers pathIntsList, reqUrl string, dirsUsers pathIntsList, reqDir string, userId int) (matchPrefix, match bool) {
	for i := range urlsUsers {
		if !util.HasUrlPrefixDir(reqUrl, urlsUsers[i].path) {
			continue
		}
		matchPrefix = true
		if userId < 0 {
			continue
		}
		for _, uid := range urlsUsers[i].values {
			if uid == userId {
				match = true
				return
			}
		}
	}

	for i := range dirsUsers {
		if !util.HasFsPrefixDir(reqDir, dirsUsers[i].path) {
			continue
		}
		matchPrefix = true
		if userId < 0 {
			continue
		}
		for _, uid := range dirsUsers[i].values {
			if uid == userId {
				match = true
				return
			}
		}
	}

	return
}

func (h *aliasHandler) getCanIndex(rawReqPath, reqFsPath string, userId int) bool {
	if h.globalIndex {
		return true
	}

	if hasUrlOrDirPrefix(h.indexUrls, rawReqPath, h.indexDirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.indexUrlsUsers, rawReqPath, h.indexDirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
}

func (h *aliasHandler) getCanUpload(info os.FileInfo, rawReqPath, reqFsPath string, userId int) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.globalUpload {
		return true
	}

	if hasUrlOrDirPrefix(h.uploadUrls, rawReqPath, h.uploadDirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.uploadUrlsUsers, rawReqPath, h.uploadDirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
}

func (h *aliasHandler) getCanMkdir(info os.FileInfo, rawReqPath, reqFsPath string, userId int) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.globalMkdir {
		return true
	}

	if hasUrlOrDirPrefix(h.mkdirUrls, rawReqPath, h.mkdirDirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.mkdirUrlsUsers, rawReqPath, h.mkdirDirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
}

func (h *aliasHandler) getCanDelete(info os.FileInfo, rawReqPath, reqFsPath string, userId int) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.globalDelete {
		return true
	}

	if hasUrlOrDirPrefix(h.deleteUrls, rawReqPath, h.deleteDirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.deleteUrlsUsers, rawReqPath, h.deleteDirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
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
