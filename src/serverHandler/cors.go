package serverHandler

import (
	"../util"
	"net/http"
	"strings"
)

func (h *handler) cors(w http.ResponseWriter, r *http.Request) {
	header := w.Header()

	header.Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodOptions {
		return
	}

	// Access-Control-Allow-Methods
	acAllowMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodTrace,
	}
	acReqMethods := r.Header["Access-Control-Request-Method"]
	if len(acReqMethods) > 0 {
		acReqMethod := acReqMethods[0]
		if !util.Contains(acAllowMethods, acReqMethod) {
			acAllowMethods = append(acAllowMethods, acReqMethod)
		}
	}
	header.Set("Access-Control-Allow-Methods", strings.Join(acAllowMethods, ", "))

	// Access-Control-Allow-Headers
	acAllowHeaders := r.Header["Access-Control-Request-Headers"]
	if len(acAllowHeaders) > 0 {
		header.Set("Access-Control-Allow-Headers", strings.Join(acAllowHeaders, ", "))
	}
}
