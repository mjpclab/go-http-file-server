package serverHandler

type alias struct {
	urlPath string
	fsPath  string
}

type aliases []*alias

func (aliases aliases) byUrlPath(urlPath string) (alias *alias, ok bool) {
	for _, alias := range aliases {
		if urlPath == alias.urlPath {
			return alias, true
		}
	}
	return nil, false
}
