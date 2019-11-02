package serverHandler

func needResponseBody(method string) bool {
	return method != "HEAD" && method != "OPTIONS" && method != "CONNECT" && method != "TRACE"
}
