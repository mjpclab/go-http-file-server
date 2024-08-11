package serverHandler

type pathContext struct {
	download     bool
	downloadfile bool
	sort         *string // keep different for param is not specified or is empty
	defaultSort  string
}

func (ctx pathContext) QueryString() string {
	// ?downloadfile&sort=x/&
	buffer := make([]byte, 1, 22)
	buffer[0] = '?' // 1 byte

	switch {
	case ctx.downloadfile:
		buffer = append(buffer, []byte("downloadfile&")...) // 13 bytes
	case ctx.download:
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

func (ctx pathContext) QueryStringOfSort(sort string) string {
	copiedCtx := ctx
	copiedCtx.sort = &sort
	return copiedCtx.QueryString()
}

func (ctx pathContext) FileQueryString() string {
	if ctx.downloadfile {
		return "?downloadfile"
	}

	return ""
}
