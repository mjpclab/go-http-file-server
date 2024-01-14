package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
	"net/http"
)

func newPathHeaders(pathHeadersEntries [][]string) pathHeadersList {
	results := make(pathHeadersList, 0, len(pathHeadersEntries))

	for _, pathHeadersSeq := range pathHeadersEntries {
		if len(pathHeadersSeq) <= 1 { // no headers
			continue
		}
		refPath := pathHeadersSeq[0]

		pathHeadersSeq = pathHeadersSeq[1:]
		headerPairCount := len(pathHeadersSeq) / 2
		if headerPairCount == 0 {
			continue
		}
		headers := make([][2]string, headerPairCount)
		for i := 0; i < headerPairCount; i++ {
			headers[i] = [2]string{pathHeadersSeq[i*2], pathHeadersSeq[i*2+1]}
		}

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
		header.Set(headers[i][0], headers[i][1])
	}
}
