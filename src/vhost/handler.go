package vhost

import (
	"../param"
	"../serverError"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"../user"
	"net/http"
)

func NewHandler(
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

	muxHandler := serverHandler.NewMultiplexer(p, *users, theme, logger)
	pathTransformHandler := serverHandler.NewPathTransformer(p.PrefixUrls, muxHandler)

	return pathTransformHandler, nil
}
