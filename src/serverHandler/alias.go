package serverHandler

type alias interface {
	urlPath() string
	fsPath() string
	isMatch(rawReqPath string) bool
	isSuccessorOf(rawReqPath string) bool
	namesEqual(a, b string) bool
	getSubPart(rawReqPath string) (subName string, isLastPart, ok bool)
}

type aliases []alias

func NewAliases(capacity int) aliases {
	aliases := make(aliases, 0, capacity)
	return aliases
}

func (aliases aliases) byUrlPath(urlPath string) (alias alias, ok bool) {
	for _, alias := range aliases {
		if alias.isMatch(urlPath) {
			return alias, true
		}
	}
	return nil, false
}
