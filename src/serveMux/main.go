package serveMux

import (
	"../param"
	"../reverseProxy"
	"../serverErrHandler"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"../user"
	"crypto/tls"
	"net/http"
	"net/url"
)

func NewServeMux(
	p *param.Param,
	logger *serverLog.Logger,
	errorHandler *serverErrHandler.ErrHandler,
) *http.ServeMux {
	users := user.NewUsers()
	for _, u := range p.UsersPlain {
		errorHandler.LogError(users.AddPlain(u.Username, u.Password))
	}
	for _, u := range p.UsersBase64 {
		errorHandler.LogError(users.AddBase64(u.Username, u.Password))
	}
	for _, u := range p.UsersMd5 {
		errorHandler.LogError(users.AddMd5(u.Username, u.Password))
	}
	for _, u := range p.UsersSha1 {
		errorHandler.LogError(users.AddSha1(u.Username, u.Password))
	}
	for _, u := range p.UsersSha256 {
		errorHandler.LogError(users.AddSha256(u.Username, u.Password))
	}
	for _, u := range p.UsersSha512 {
		errorHandler.LogError(users.AddSha512(u.Username, u.Password))
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: p.IgnoreProxyTargetBadCert},
	}
	fallbackProxies := mapToReverseProxy(p.FallbackProxies, tr)
	alwaysProxies := mapToReverseProxy(p.AlwaysProxies, tr)

	tplObj, err := tpl.LoadPage(p.Template)
	errorHandler.LogError(err)

	aliases := p.Aliases
	if _, hasRootAlias := aliases["/"]; !hasRootAlias {
		aliases["/"] = p.Root
	}

	handlers := map[string]http.Handler{}
	for urlPath, fsPath := range aliases {
		handlers[urlPath] = serverHandler.NewHandler(fsPath, urlPath, p, users, fallbackProxies, alwaysProxies, tplObj, logger, errorHandler)
	}

	// create ServeMux
	serveMux := &http.ServeMux{}
	for urlPath, handler := range handlers {
		serveMux.Handle(urlPath, handler)
		if len(urlPath) > 1 {
			serveMux.Handle(urlPath+"/", handler)
		}
	}

	return serveMux
}

func mapToReverseProxy(input map[string]string, tr http.RoundTripper) map[string]http.Handler {
	maps := map[string]http.Handler{}

	for inUrl, target := range input {
		targetUrl, err := url.Parse(target)
		if err != nil || len(targetUrl.Scheme) == 0 || len(targetUrl.Host) == 0 {
			continue
		}
		var proxyHandler http.Handler = reverseProxy.NewReverseProxy(targetUrl, tr)
		proxyHandler = http.StripPrefix(inUrl, proxyHandler)
		maps[inUrl] = proxyHandler
	}

	return maps
}
