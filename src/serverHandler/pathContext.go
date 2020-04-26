package serverHandler

type pathContext struct {
	sort string
}

func (ctx *pathContext) QueryString() string {
	if len(ctx.sort) > 0 {
		return "?sort=" + ctx.sort
	} else {
		return ""
	}
}

func (ctx *pathContext) QueryStringOfSort(sort string) string {
	copiedCtx := *ctx
	copiedCtx.sort = sort
	return copiedCtx.QueryString()
}
