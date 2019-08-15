package param

import "flag"

type Param struct {
	Root     string
	Key      string
	Cert     string
	Listen   string
	Template string
}

var param Param

func init() {
	flag.StringVar(&param.Root, "root", ".", "root directory of server")
	flag.StringVar(&param.Key, "key", "", "TLS certificate key path")
	flag.StringVar(&param.Cert, "cert", "", "TLS certificate path")
	flag.StringVar(&param.Listen, "listen", "", "address and port to listen")
	flag.StringVar(&param.Template, "template", "", "page template path")

	flag.Parse()
}

func Parse() *Param {
	paramCopied := param
	return &paramCopied
}
