package serverHandler

import "mjpclab.dev/ghfs/src/util"

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

func (ha *hierarchyAvailability) match(urlPath, fsPath string, userId int) bool {
	if ha.global {
		return true
	}

	if hasUrlOrDirPrefix(ha.urls, urlPath, ha.dirs, fsPath) {
		return true
	}

	if userId >= 0 {
		if _, match := hasUrlOrDirPrefixUsers(ha.urlsUsers, urlPath, ha.dirsUsers, fsPath, userId); match {
			return true
		}
	}

	return false
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
