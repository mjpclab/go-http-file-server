package serverHandler

import (
	"../param"
	"../serverError"
	"../serverLog"
	"../tpl"
	"../user"
	"net/http"
	"strings"
)

func NewVhostHandler(
	p *param.Param,
	logger *serverLog.Logger,
	theme tpl.Theme,
) (handler http.Handler, errs []error) {
	// users
	users := user.NewList(p.UserMatchCase)
	for _, u := range p.UsersPlain {
		errs = serverError.AppendError(errs, users.AddPlain(u.Username, u.Password))
	}
	for _, u := range p.UsersBase64 {
		errs = serverError.AppendError(errs, users.AddBase64(u.Username, u.Password))
	}
	for _, u := range p.UsersMd5 {
		errs = serverError.AppendError(errs, users.AddMd5(u.Username, u.Password))
	}
	for _, u := range p.UsersSha1 {
		errs = serverError.AppendError(errs, users.AddSha1(u.Username, u.Password))
	}
	for _, u := range p.UsersSha256 {
		errs = serverError.AppendError(errs, users.AddSha256(u.Username, u.Password))
	}
	for _, u := range p.UsersSha512 {
		errs = serverError.AppendError(errs, users.AddSha512(u.Username, u.Password))
	}

	if len(errs) > 0 {
		return nil, errs
	}

	restrictAccessUrls := newRestrictAccesses(p.RestrictAccessUrls)
	restrictAccessDirs := newRestrictAccesses(p.RestrictAccessDirs)
	restrictAccess := hasRestrictAccess(p.GlobalRestrictAccess, restrictAccessUrls, restrictAccessDirs)
	pageVaryV1 := "Accept-Encoding"
	contentVaryV1 := ""
	if restrictAccess {
		pageVaryV1 += ", Referer"
		contentVaryV1 = "Referer"
	}
	pageVary := strings.ToLower(pageVaryV1)
	contentVary := strings.ToLower(contentVaryV1)

	ap := &aliasParam{
		users:  *users,
		theme:  theme,
		logger: logger,

		headersUrls: newPathHeaders(p.HeadersUrls),
		headersDirs: newPathHeaders(p.HeadersDirs),

		restrictAccess:     restrictAccess,
		restrictAccessUrls: restrictAccessUrls,
		restrictAccessDirs: restrictAccessDirs,

		pageVaryV1:    pageVaryV1,
		pageVary:      pageVary,
		contentVaryV1: contentVaryV1,
		contentVary:   contentVary,
	}

	muxHandler := newMultiplexHandler(p, ap)
	pathTransformHandler := newPathTransformHandler(p.PrefixUrls, muxHandler)

	return pathTransformHandler, nil
}
