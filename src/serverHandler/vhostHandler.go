package serverHandler

import (
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverError"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/tpl/defaultTheme"
	"mjpclab.dev/ghfs/src/tpl/theme"
	"mjpclab.dev/ghfs/src/user"
	"net/http"
	"regexp"
)

type vhostContext struct {
	logger *serverLog.Logger
	theme  theme.Theme

	users            *user.List
	authUrlsUsers    pathIntsList
	authDirsUsers    pathIntsList
	indexUrlsUsers   pathIntsList
	indexDirsUsers   pathIntsList
	uploadUrlsUsers  pathIntsList
	uploadDirsUsers  pathIntsList
	mkdirUrlsUsers   pathIntsList
	mkdirDirsUsers   pathIntsList
	deleteUrlsUsers  pathIntsList
	deleteDirsUsers  pathIntsList
	archiveUrlsUsers pathIntsList
	archiveDirsUsers pathIntsList

	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp

	restrictAccessUrls pathStringsList
	restrictAccessDirs pathStringsList

	headersUrls pathHeadersList
	headersDirs pathHeadersList
}

func NewVhostHandler(
	p *param.Param,
	logger *serverLog.Logger,
	theme theme.Theme,
) (handler http.Handler, errs []error) {
	// users
	users := user.NewList()
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

	if logger == nil {
		logger = &serverLog.Logger{}
	}

	if theme == nil {
		theme = defaultTheme.DefaultTheme
	}

	// alias param
	vhostCtx := &vhostContext{
		logger: logger,
		theme:  theme,

		users:            users,
		authUrlsUsers:    pathUsernamesToPathUids(users, p.AuthUrlsUsers),
		authDirsUsers:    pathUsernamesToPathUids(users, p.AuthDirsUsers),
		indexUrlsUsers:   pathUsernamesToPathUids(users, p.IndexUrlsUsers),
		indexDirsUsers:   pathUsernamesToPathUids(users, p.IndexDirsUsers),
		uploadUrlsUsers:  pathUsernamesToPathUids(users, p.UploadUrlsUsers),
		uploadDirsUsers:  pathUsernamesToPathUids(users, p.UploadDirsUsers),
		mkdirUrlsUsers:   pathUsernamesToPathUids(users, p.MkdirUrlsUsers),
		mkdirDirsUsers:   pathUsernamesToPathUids(users, p.MkdirDirsUsers),
		deleteUrlsUsers:  pathUsernamesToPathUids(users, p.DeleteUrlsUsers),
		deleteDirsUsers:  pathUsernamesToPathUids(users, p.DeleteDirsUsers),
		archiveUrlsUsers: pathUsernamesToPathUids(users, p.ArchiveUrlsUsers),
		archiveDirsUsers: pathUsernamesToPathUids(users, p.ArchiveDirsUsers),

		shows:     shows,
		showDirs:  showDirs,
		showFiles: showFiles,
		hides:     hides,
		hideDirs:  hideDirs,
		hideFiles: hideFiles,

		restrictAccessUrls: newRestrictAccesses(p.RestrictAccessUrls),
		restrictAccessDirs: newRestrictAccesses(p.RestrictAccessDirs),

		headersUrls: newPathHeaders(p.HeadersUrls),
		headersDirs: newPathHeaders(p.HeadersDirs),
	}

	handler = newMultiplexHandler(p, vhostCtx)
	handler = newPreprocessHandler(logger, p.PreMiddlewares, handler)
	handler = newPathTransformHandler(p.PrefixUrls, handler)
	return
}
