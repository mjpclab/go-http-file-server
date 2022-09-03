package middleware

import "net/http"

type ProcessResult int

const (
	GoNext ProcessResult = iota
	SkipRests
	Processed
)

type Middleware func(w http.ResponseWriter, r *http.Request, context *Context) (result ProcessResult)
