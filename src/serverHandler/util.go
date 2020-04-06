package serverHandler

import (
	"net/http"
	"os"
)

func needResponseBody(method string) bool {
	return method != http.MethodHead &&
		method != http.MethodOptions &&
		method != http.MethodConnect &&
		method != http.MethodTrace
}

func containsItem(infos []os.FileInfo, name string) bool {
	for i := range infos {
		if infos[i].Name() == name {
			return true
		}
	}
	return false
}
