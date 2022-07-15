package serverHandler

import (
	"net/http"
)

func cors(w http.ResponseWriter) {
	header := w.Header()

	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "*")
	header.Set("Access-Control-Allow-Headers", "*")
}
