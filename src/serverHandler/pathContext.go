package serverHandler

type pathContext struct {
	simple      bool
	download    bool
	sort        *string // keep different for param is not specified or is empty
	defaultSort string
}

func (ctx pathContext) QueryString() string {
	// ?simpledownload&sort=x/&
	buffer := make([]byte, 1, 24)
	buffer[0] = '?' // 1 byte

	switch {
	case ctx.simple && ctx.download:
		buffer = append(buffer, []byte("simpledownload&")...) // 15 bytes
	case ctx.simple:
		buffer = append(buffer, []byte("simple&")...) // 7 bytes
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

func (ctx pathContext) SubFileQueryString() string {
	if ctx.download {
		return "?download"
	}

	return ""
}
