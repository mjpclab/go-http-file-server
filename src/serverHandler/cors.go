package serverHandler

import (
	"../shimgo"
	"../util"
	"net/http"
	"strings"
)

func (h *handler) cors(w http.ResponseWriter, r *http.Request) {
	header := w.Header()

	header.Set("Access-Control-Allow-Origin", "*")

	if r.Method != shimgo.Net_Http_MethodOptions {
		return
	}

	// Access-Control-Allow-Methods
	acAllowMethods := []string{
		shimgo.Net_Http_MethodGet,
		shimgo.Net_Http_MethodHead,
		shimgo.Net_Http_MethodPost,
		shimgo.Net_Http_MethodPut,
		shimgo.Net_Http_MethodDelete,
		shimgo.Net_Http_MethodOptions,
		shimgo.Net_Http_MethodTrace,
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
