package serverHandler

type pathContext struct {
	defaultSort string
	sort        *string
}

func (ctx *pathContext) QueryString() string {
	if ctx.sort != nil && *(ctx.sort) != ctx.defaultSort {
		return "?sort=" + *(ctx.sort)
	} else {
		return ""
	}
}

func (ctx *pathContext) QueryStringOfSort(sort string) string {
	copiedCtx := *ctx
	copiedCtx.sort = &sort
	return copiedCtx.QueryString()
}
