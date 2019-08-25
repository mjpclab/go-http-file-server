package param

import (
	argParser "../goNixArgParser"
	"../util"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)
import "../serverError"

type Param struct {
	Root     string
	Aliases  map[string]string
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

	err = argParser.AddFlagsValues("aliases", []string{"-a", "--alias"}, nil, "set alias path, <sep><url><sep><path>, e.g. :/doc:/usr/share/doc")
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

	// normalize aliases
	param.Aliases = map[string]string{}
	arrAlias := result.GetValues("aliases")
	if arrAlias != nil {
		for _, alias := range arrAlias {
			sep, sepLen := utf8.DecodeRuneInString(alias)
			if sepLen == 0 {
				continue
			}
			alias = alias[sepLen:]
			sepIndex := strings.Index(alias, string(sep))

			urlPath := util.CleanUrlPath(alias[:sepIndex])
			fsPath := path.Clean(alias[sepIndex+sepLen:])
			if len(fsPath) == 0 {
				fsPath = "."
			}

			param.Aliases[urlPath] = fsPath
		}
	}
}

func Parse() *Param {
	paramCopied := param
	return &paramCopied
}
