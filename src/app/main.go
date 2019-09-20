package app

import (
	"../param"
	"../serverErrHandler"
	"../vhost"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

type App struct {
	vhosts  []*vhost.VHost
	listens Listens
}

func (app *App) Open() {
	var err error
	hasError := false

	// create listener
	for _, l := range app.listens {
		if l.proto == "unix" {
			sockInfo, _ := os.Lstat(l.addr)
			if sockInfo != nil && (sockInfo.Mode()&os.ModeSocket != 0) {
				os.Remove(l.addr)
			}
		}

		l.listener, err = net.Listen(l.proto, l.addr)
		if err != nil {
			hasError = true
			serverErrHandler.CheckError(err)
		}

		if err == nil && l.proto == "unix" {
			os.Chmod(l.addr, 0777)
		}
	}

	if hasError {
		return
	}

	// create server
	wgStart := sync.WaitGroup{}
	wgStop := sync.WaitGroup{}
	for _, l := range app.listens {
		item := l

		server := &http.Server{
			Handler: http.HandlerFunc(item.handlerFunc),
		}
		if item.useTLS {
			tlsConfig := &tls.Config{
				Certificates: item.certs,
			}
			tlsConfig.BuildNameToCertificate()

			server.TLSConfig = tlsConfig
		}

		wgStart.Add(1)
		wgStop.Add(1)
		go func() {
			wgStart.Done()
			if item.useTLS {
				err = server.ServeTLS(item.listener, "", "")
			} else {
				err = server.Serve(item.listener)
			}
			if err != nil {
				hasError = true
				serverErrHandler.CheckError(err)
			}
			wgStop.Done()
		}()
	}

	wgStart.Wait()
	if hasError {
		app.Close()
		return
	}
	wgStop.Wait()
}

func (app *App) Close() {
	for _, vh := range app.vhosts {
		vh.Close()
	}

	for _, l := range app.listens {
		if l.server != nil {
			l.server.Close()
			l.server = nil
		}

		if l.listener != nil {
			l.listener.Close()
			l.listener = nil
		}
	}
}

func (app *App) ReOpenLog() {
	for _, vh := range app.vhosts {
		vh.ReOpenLog()
	}
}

func NewApp(params []*param.Param) *App {
	app := &App{
		listens: Listens{},
	}

	// vhosts
	vhosts := []*vhost.VHost{}
	for _, p := range params {
		vh := vhost.NewVHost(p)
		vhosts = append(vhosts, vh)
	}
	app.vhosts = vhosts

	// listens
	for _, vh := range vhosts {
		hasErr := false

		// verify
		for _, vhListen := range vh.Listens {
			vhAddr := vhListen.Addr

			// listen -> useTLS conflicts
			item := app.listens.findItemByAddr(vhAddr)
			if item != nil && item.useTLS != vhListen.UseTLS {
				hasErr = true
				serverErrHandler.CheckError(errors.New(vhAddr + " cannot served for both PLAIN and TLS mode"))
			}

			// listen, hostname duplicated
			for _, vhHostname := range vh.Hostnames {
				item := app.listens.findItemByAddrHostname(vhAddr, vhHostname)
				if item != nil {
					hasErr = true
					serverErrHandler.CheckError(errors.New(vhAddr + " " + vhHostname + " duplicated Listen and Hostname"))
				}
			}
		}

		// create or update ListenItem
		for _, vhListen := range vh.Listens {
			// construct ListenItem if not exists
			item := app.listens.findItemByAddr(vhListen.Addr)
			if item == nil {
				item = &ListenItem{
					proto:     vhListen.Proto,
					addr:      vhListen.Addr,
					useTLS:    vhListen.UseTLS,
					hostnames: []string{},
					certs:     []tls.Certificate{},
					server:    nil,
					vhosts:    []*vhost.VHost{},
				}
				item.handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					hostname := r.Host
					colonIndex := strings.LastIndexByte(hostname, ':')
					if colonIndex >= 0 {
						hostname = hostname[:colonIndex]
					}

					var serveVHost *vhost.VHost
					for _, vh := range item.vhosts {
						if vh.MatchHostname(hostname) {
							serveVHost = vh
							break
						}
					}
					if serveVHost == nil {
						serveVHost = item.vhosts[0]
					}

					serveVHost.Mux.ServeHTTP(w, r)
				}

				app.listens = append(app.listens, item)
			}

			// update
			item.hostnames = append(item.hostnames, vh.Hostnames...)
			if item.useTLS {
				cert, err := tls.LoadX509KeyPair(vhListen.Cert, vhListen.Key)
				if err != nil {
					hasErr = true
					serverErrHandler.CheckError(err)
				}
				item.certs = append(item.certs, cert)
			}
			item.vhosts = append(item.vhosts, vh)

			// verified
			app.listens = append(app.listens, )
		}

		if hasErr {
			os.Exit(1)
		}
	}

	return app
}
