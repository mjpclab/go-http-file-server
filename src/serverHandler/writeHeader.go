package serverHandler

import "net/http"

func writeHeader(w http.ResponseWriter, r *http.Request, data *responseData) {
	switch {
	case data.HasForbiddenError:
		w.WriteHeader(http.StatusForbidden)
	case data.HasNotFoundError:
		w.WriteHeader(http.StatusNotFound)
	case data.HasInternalError:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
