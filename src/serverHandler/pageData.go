package serverHandler

import (
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
	Item     os.FileInfo
	SubItems []os.FileInfo
	Error    error
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https:"
	} else {
		return "http:"
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

func readdir(realPath string) (item os.FileInfo, subItems []os.FileInfo, err error) {
	var f *os.File
	f, err = os.Open(realPath)
	if err != nil {
		return
	}
	defer f.Close()

	item, err = f.Stat()
	if err != nil {
		return
	}

	if !item.IsDir() {
		return
	}

	subItems, err = f.Readdir(-1)
	if err == nil {
		sortSubItems(subItems)
	}

	return
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

func getPageData(root string, r *http.Request) *pageData {
	requestPath := path.Clean(r.URL.Path)
	realPath := path.Clean(root + requestPath)

	scheme := getScheme(r)
	relPath := requestPath[1:]
	pathEntries := getPathEntries(relPath)
	item, subItems, err := readdir(realPath)

	data := &pageData{
		Scheme:   scheme,
		Host:     r.Host,
		Path:     relPath,
		Paths:    pathEntries,
		Item:     item,
		SubItems: subItems,
		Error:    err,
	}

	return data
}
