package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
	"os"
)

type hierarchyAvailability struct {
	global    bool
	urls      []string
	urlsUsers pathIntsList
	dirs      []string
	dirsUsers pathIntsList
}

func newHierarchyAvailability(
	baseUrl, baseDir string,
	global bool,
	allUrls []string, allUrlsUsers pathIntsList,
	allDirs []string, allDirsUsers pathIntsList,
) *hierarchyAvailability {
	return &hierarchyAvailability{
		global:    global || prefixMatched(allUrls, util.HasUrlPrefixDir, baseUrl) || prefixMatched(allDirs, util.HasFsPrefixDir, baseDir),
		urls:      filterSuccessor(allUrls, util.HasUrlPrefixDir, baseUrl),
		urlsUsers: allUrlsUsers.filterSuccessor(true, util.HasUrlPrefixDir, baseUrl),
		dirs:      filterSuccessor(allDirs, util.HasFsPrefixDir, baseDir),
		dirsUsers: allDirsUsers.filterSuccessor(true, util.HasFsPrefixDir, baseDir),
	}
}

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
	if h.index.global {
		return true
	}

	if hasUrlOrDirPrefix(h.index.urls, rawReqPath, h.index.dirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.index.urlsUsers, rawReqPath, h.index.dirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
}

func (h *aliasHandler) getCanUpload(info os.FileInfo, rawReqPath, reqFsPath string, userId int) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.upload.global {
		return true
	}

	if hasUrlOrDirPrefix(h.upload.urls, rawReqPath, h.upload.dirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.upload.urlsUsers, rawReqPath, h.upload.dirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
}

func (h *aliasHandler) getCanMkdir(info os.FileInfo, rawReqPath, reqFsPath string, userId int) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.mkdir.global {
		return true
	}

	if hasUrlOrDirPrefix(h.mkdir.urls, rawReqPath, h.mkdir.dirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.mkdir.urlsUsers, rawReqPath, h.mkdir.dirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
}

func (h *aliasHandler) getCanDelete(info os.FileInfo, rawReqPath, reqFsPath string, userId int) bool {
	if info == nil || !info.IsDir() {
		return false
	}

	if h.delete.global {
		return true
	}

	if hasUrlOrDirPrefix(h.delete.urls, rawReqPath, h.delete.dirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.delete.urlsUsers, rawReqPath, h.delete.dirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
}

func (h *aliasHandler) getCanArchive(subInfos []os.FileInfo, rawReqPath, reqFsPath string, userId int) bool {
	if len(subInfos) == 0 {
		return false
	}

	if h.archive.global {
		return true
	}

	if hasUrlOrDirPrefix(h.archive.urls, rawReqPath, h.archive.dirs, reqFsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(h.archive.urlsUsers, rawReqPath, h.archive.dirsUsers, reqFsPath, userId); match {
			return true
		}
	}

	return false
}

func (h *aliasHandler) getCanCors(rawReqPath, reqFsPath string) bool {
	if h.cors.global {
		return true
	}

	return hasUrlOrDirPrefix(h.cors.urls, rawReqPath, h.cors.dirs, reqFsPath)
}
