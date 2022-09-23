package serverHandler

import (
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverError"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/tpl"
	"mjpclab.dev/ghfs/src/user"
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

	// show/hide
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

	// restrict access
	restrictAccessUrls := newRestrictAccesses(p.RestrictAccessUrls)
	restrictAccessDirs := newRestrictAccesses(p.RestrictAccessDirs)
	restrictAccess := hasRestrictAccess(p.GlobalRestrictAccess, restrictAccessUrls, restrictAccessDirs)

	// `Vary` header
	pageVarys := make([]string, 0, 4)
	contentVarys := make([]string, 0, 3)
	pageVarys = append(pageVarys, "Accept-Encoding")
	if restrictAccess {
		pageVarys = append(pageVarys, "Referer", "Origin")
		contentVarys = append(contentVarys, "Referer", "Origin")
	}
	if p.GlobalAuth || len(p.AuthUrls) > 0 || len(p.AuthDirs) > 0 {
		pageVarys = append(pageVarys, "Authorization")
		contentVarys = append(contentVarys, "Authorization")
	}

	pageVaryV1 := strings.Join(pageVarys, ", ")
	contentVaryV1 := strings.Join(contentVarys, ", ")
	pageVary := strings.ToLower(pageVaryV1)
	contentVary := strings.ToLower(contentVaryV1)

	// alias param
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
	preprocessHandler := newPreprocessHandler(logger, p.PreMiddlewares, muxHandler)
	pathTransformHandler := newPathTransformHandler(p.PrefixUrls, preprocessHandler)

	return pathTransformHandler, nil
}
