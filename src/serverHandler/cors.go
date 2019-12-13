package serverHandler

import (
	"../util"
	"net/http"
	"strings"
)

func (h *handler) cors(w http.ResponseWriter, r *http.Request) {
	header := w.Header()

	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Credentials", "true")

	if r.Method != "OPTIONS" {
		return
	}

	// Access-Control-Allow-Methods
	acAllowMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "TRACE"}
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
