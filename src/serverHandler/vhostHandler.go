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
		errs = serverError.AppendError(errs, users.AddPlain(u[0], u[1]))
	}
	for _, u := range p.UsersBase64 {
		errs = serverError.AppendError(errs, users.AddBase64(u[0], u[1]))
	}
	for _, u := range p.UsersMd5 {
		errs = serverError.AppendError(errs, users.AddMd5(u[0], u[1]))
	}
	for _, u := range p.UsersSha1 {
		errs = serverError.AppendError(errs, users.AddSha1(u[0], u[1]))
	}
	for _, u := range p.UsersSha256 {
		errs = serverError.AppendError(errs, users.AddSha256(u[0], u[1]))
	}
	for _, u := range p.UsersSha512 {
		errs = serverError.AppendError(errs, users.AddSha512(u[0], u[1]))
	}

	shows, err := wildcardToRegexp(p.Shows)
	errs = serverError.AppendError(errs, err)
	showDirs, err := wildcardToRegexp(p.ShowDirs)
	errs = serverError.AppendError(errs, err)
	showFiles, err := wildcardToRegexp(p.ShowFiles)
	errs = serverError.AppendError(errs, err)
	hides, err := wildcardToRegexp(p.Hides)
	errs = serverError.AppendError(errs, err)
	hideDirs, err := wildcardToRegexp(p.HideDirs)
	errs = serverError.AppendError(errs, err)
	hideFiles, err := wildcardToRegexp(p.HideFiles)
	errs = serverError.AppendError(errs, err)

	if len(errs) > 0 {
		return nil, errs
	}

	restrictAccessUrls := newRestrictAccesses(p.RestrictAccessUrls)
	restrictAccessDirs := newRestrictAccesses(p.RestrictAccessDirs)
	restrictAccess := hasRestrictAccess(p.GlobalRestrictAccess, restrictAccessUrls, restrictAccessDirs)
	pageVaryV1 := "Accept-Encoding"
	contentVaryV1 := ""
	if restrictAccess {
		pageVaryV1 += ", Referer, Origin"
		contentVaryV1 = "Referer, Origin"
	}
	pageVary := strings.ToLower(pageVaryV1)
	contentVary := strings.ToLower(contentVaryV1)

	ap := &aliasParam{
		users:  *users,
		theme:  theme,
		logger: logger,

		shows:     shows,
		showDirs:  showDirs,
		showFiles: showFiles,
		hides:     hides,
		hideDirs:  hideDirs,
		hideFiles: hideFiles,

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
