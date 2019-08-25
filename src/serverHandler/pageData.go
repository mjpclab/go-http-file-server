package serverHandler

import (
	"../util"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

type pathEntry struct {
	Name string
	Path string
}

type pageData struct {
	Scheme   string
	Host     string
	Path     string
	Paths    []*pathEntry
	File     *os.File
	Item     os.FileInfo
	SubItems []os.FileInfo
	Errors   []error
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https:"
	} else {
		return "http:"
	}
}

func (h *handler) readdir(rawRequestPath, requestPath string, errors *[]error) (file *os.File, item os.FileInfo, subItems []os.FileInfo) {
	var err error

	fsPath := path.Clean(h.root + requestPath)

	file, err = os.Open(fsPath)
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	item, err = file.Stat()
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	if !item.IsDir() {
		return
	}

	subItems, err = file.Readdir(-1)
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	return
}

func (h *handler) mergeAlias(rawRequestPath string, subItems *[]os.FileInfo, errors *[]error) {
	for aliasUrlPath, aliasFsPath := range h.aliases {
		if len(aliasUrlPath) <= len(rawRequestPath) {
			continue
		}

		suffixIndex := len(rawRequestPath)
		aliasPrefix := aliasUrlPath[:suffixIndex]
		aliasSuffix := aliasUrlPath[suffixIndex:]

		if aliasPrefix != rawRequestPath {
			continue
		}

		if len(aliasPrefix) != 1 && aliasSuffix[0] != '/' {
			// aliasUrlPath doesn't cover the whole directory name
			// e.g:
			// rawRequestPath == "/abc/def/ghi"
			// aliasPrefix == "/abc/de"
			continue
		}
		if aliasSuffix[0] == '/' {
			aliasSuffix = aliasSuffix[1:]
		}

		slashIndex := strings.Index(aliasSuffix, "/")
		var nextName string
		if slashIndex >= 0 {
			nextName = aliasSuffix[:slashIndex]
		} else {
			nextName = aliasSuffix
		}

		var aliasSubItem os.FileInfo
		if path.Dir(aliasUrlPath) == rawRequestPath { // reached second deepest path of alias
			var err error
			aliasSubItem, err = os.Stat(aliasFsPath)
			if err == nil {
				aliasSubItem = NewRenamedFileInfo(nextName, aliasSubItem)
			} else {
				*errors = append(*errors, err)
				aliasSubItem = NewFakeFileInfo(nextName, true)
			}
		} else {
			aliasSubItem = NewFakeFileInfo(nextName, true)
		}

		replaced := false
		for i, subItem := range *subItems {
			if subItem.Name() == nextName {
				(*subItems)[i] = aliasSubItem
				replaced = true
				break
			}
		}

		if !replaced {
			*subItems = append(*subItems, aliasSubItem)
		}
	}
}

func sortSubItems(subItems []os.FileInfo) {
	sort.Slice(
		subItems,
		func(i, j int) bool {
			itemI := subItems[i]
			itemJ := subItems[j]

			isDirI := itemI.IsDir()
			isDirJ := itemJ.IsDir()

			if (isDirI && isDirJ) || (!isDirI && !isDirJ) {
				return itemI.Name() < itemJ.Name()
			}

			return isDirI
		},
	)
}

func getPathEntries(path string) []*pathEntry {
	var pathParts []string
	if len(path) > 0 {
		pathParts = strings.Split(path, "/")
	} else {
		pathParts = []string{}
	}

	pathEntries := make([]*pathEntry, len(pathParts))
	for i, part := range pathParts {
		pathEntries = append(pathEntries, &pathEntry{
			Name: part,
			Path: "/" + strings.Join(pathParts[:i+1], "/"),
		})
	}

	return pathEntries
}

func (h *handler) getPageData(r *http.Request) *pageData {
	rawRequestPath := util.CleanUrlPath(r.URL.Path)
	requestPath := util.CleanUrlPath(rawRequestPath[len(h.urlPrefix):]) // strip url prefix path

	scheme := getScheme(r)

	relPath := rawRequestPath[1:]
	pathEntries := getPathEntries(relPath)

	errors := []error{}
	file, item, subItems := h.readdir(rawRequestPath, requestPath, &errors)
	h.mergeAlias(rawRequestPath, &subItems, &errors)
	sortSubItems(subItems)

	data := &pageData{
		Scheme:   scheme,
		Host:     r.Host,
		Path:     relPath,
		Paths:    pathEntries,
		File:     file,
		Item:     item,
		SubItems: subItems,
		Errors:   errors,
	}

	return data
}
