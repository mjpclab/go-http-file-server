package middleware

import "net/http"

type Middleware func(w http.ResponseWriter, r *http.Request, context *Context) (processed bool)
