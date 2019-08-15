package param

import "flag"

type Param struct {
	Root     string
	Listen   string
	Key      string
	Cert     string
	Template string
}

var param Param

func init() {
	flag.StringVar(&param.Root, "root", ".", "root directory of server")
	flag.StringVar(&param.Listen, "listen", ":80", "address and port to listen")
	flag.StringVar(&param.Key, "key", "", "TLS certificate key path")
	flag.StringVar(&param.Cert, "cert", "", "TLS certificate path")
	flag.StringVar(&param.Template, "template", "", "page template path")

	flag.Parse()
}

func Parse() *Param {
	paramCopied := param
	return &paramCopied
}
