package param

import (
	argParser "../goNixArgParser"
	"os"
)
import "../serverError"

type Param struct {
	Root     string
	Key      string
	Cert     string
	Listen   string
	Template string
}

var param Param

func init() {
	// define option
	var err error
	err = argParser.AddFlagsValue("root", []string{"-r", "--root"}, ".", "root directory of server")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("key", []string{"-k", "--key"}, "", "TLS certificate key path")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("cert", []string{"-c", "--cert"}, "", "TLS certificate path")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("listen", []string{"-l", "--listen"}, "", "address and port to listen")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("template", []string{"-t", "--template"}, "", "address and port to listen")
	serverError.CheckFatal(err)

	err = argParser.AddFlags("help", []string{"-h", "--help"}, "print this help")
	serverError.CheckFatal(err)

	// parse option
	result := argParser.Parse()

	// help
	if result.HasKey("help") {
		argParser.PrintHelp()
		os.Exit(0)
	}

	// normalize option
	param.Root = result.GetValue("root")
	param.Key = result.GetValue("key")
	param.Cert = result.GetValue("cert")
	if result.HasKey("listen") {
		param.Listen = result.GetValue("listen")
	} else {
		rests := result.GetRests()
		if len(rests) > 0 {
			param.Listen = rests[0]
		}
	}
	param.Template = result.GetValue("template")
}

func Parse() *Param {
	paramCopied := param
	return &paramCopied
}
