package serverHandler

import (
	"../util"
	"net/http"
)

type pathHeaders struct {
	path    string
	headers [][2]string
}

func newPathHeaders(pathHeadersMap map[string][][2]string) []pathHeaders {
	results := make([]pathHeaders, 0, len(pathHeadersMap))

	for refPath, headers := range pathHeadersMap {
		results = append(results, pathHeaders{refPath, headers})
	}

	return results
}

func (h *aliasHandler) getHeaders(reqUrlPath, reqFsPath string, doGetHeaders bool) [][2]string {
	if !doGetHeaders {
		return nil
	}

	headers := make([][2]string, len(h.globalHeaders), len(h.globalHeaders)+len(h.headersUrls)+len(h.headersDirs))

	if len(h.globalHeaders) > 0 {
		copy(headers, h.globalHeaders)
	}

	for i := range h.headersUrls {
		if util.HasUrlPrefixDir(reqUrlPath, h.headersUrls[i].path) {
			headers = append(headers, h.headersUrls[i].headers...)
		}
	}

	for i := range h.headersDirs {
		if util.HasFsPrefixDir(reqFsPath, h.headersDirs[i].path) {
			headers = append(headers, h.headersDirs[i].headers...)
		}
	}

	return headers
}

func header(w http.ResponseWriter, headers [][2]string) {
	header := w.Header()
	for i := range headers {
		header.Add(headers[i][0], headers[i][1])
	}
}
