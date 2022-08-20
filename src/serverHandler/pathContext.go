package serverHandler

type pathContext struct {
	download    bool
	sort        *string // keep different for param is not specified or is empty
	defaultSort string
}

func (ctx pathContext) QueryString() string {
	// ?download&sort=x/
	buffer := make([]byte, 1, 18)
	buffer[0] = '?'

	if ctx.download {
		buffer = append(buffer, []byte("download&")...) // 9 bytes
	}

	if ctx.sort != nil && *(ctx.sort) != ctx.defaultSort {
		buffer = append(buffer, []byte("sort=")...)   // 5 bytes
		buffer = append(buffer, []byte(*ctx.sort)...) // 2 bytes
		buffer = append(buffer, '&')                  // 1 byte
	}

	buffer = buffer[:len(buffer)-1]
	return string(buffer)
}

func (ctx pathContext) FileQueryString() string {
	if ctx.download {
		return "?download"
	}

	return ""
}

func (ctx pathContext) QueryStringOfSort(sort string) string {
	copiedCtx := ctx
	copiedCtx.sort = &sort
	return copiedCtx.QueryString()
}
