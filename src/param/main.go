package param

import (
	argParser "../goNixArgParser"
	"../util"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"
)
import "../serverError"

type Param struct {
	Root      string
	Aliases   map[string]string
	Uploads   map[string]bool
	Key       string
	Cert      string
	Listen    string
	Template  string
	Shows     *regexp.Regexp
	ShowDirs  *regexp.Regexp
	ShowFiles *regexp.Regexp
	Hides     *regexp.Regexp
	HideDirs  *regexp.Regexp
	HideFiles *regexp.Regexp
	AccessLog string
	ErrorLog  string
}

var param Param

func getWildcardRegexp(whildcards []string) (*regexp.Regexp, error) {
	if len(whildcards) > 0 {
		for i, show := range whildcards {
			whildcards[i] = util.WildcardToRegexp(show)
		}
		exp := strings.Join(whildcards, "|")
		return regexp.Compile(exp)
	}

	return nil, nil
}

func init() {
	// define option
	var err error

	err = argParser.AddFlagsValue("root", []string{"-r", "--root"}, ".", "root directory of server")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValues("aliases", []string{"-a", "--alias"}, nil, "set alias path, <sep><url><sep><path>, e.g. :/doc:/usr/share/doc")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValues("uploads", []string{"-u", "--upload"}, nil, "url path that allow upload files")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("key", []string{"-k", "--key"}, "", "TLS certificate key path")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("cert", []string{"-c", "--cert"}, "", "TLS certificate path")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("listen", []string{"-l", "--listen"}, "", "address and port to listen")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("template", []string{"-t", "--template"}, "", "custom template file for page")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValues("shows", []string{"-S", "--show"}, nil, "show directories or files match wildcard")
	serverError.CheckFatal(err)
	err = argParser.AddFlagsValues("showdirs", []string{"-SD", "--show-dir"}, nil, "show directories match wildcard")
	serverError.CheckFatal(err)
	err = argParser.AddFlagsValues("showfiles", []string{"-SF", "--show-file"}, nil, "show files match wildcard")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValues("hides", []string{"-H", "--hide"}, nil, "hide directories or files match wildcard")
	serverError.CheckFatal(err)
	err = argParser.AddFlagsValues("hidedirs", []string{"-HD", "--hide-dir"}, nil, "hide directories match wildcard")
	serverError.CheckFatal(err)
	err = argParser.AddFlagsValues("hidefiles", []string{"-HF", "--hide-file"}, nil, "hide files match wildcard")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("accesslog", []string{"-L", "--access-log"}, "", "access log file, use \"-\" for stdout")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("errorlog", []string{"-E", "--error-log"}, "-", "error log file, use \"-\" for stderr")
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
	param.AccessLog = result.GetValue("accesslog")
	param.ErrorLog = result.GetValue("errorlog")

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

	// normalize uploads
	param.Uploads = map[string]bool{}
	arrUploads := result.GetValues("uploads")
	if arrUploads != nil {
		for _, upload := range arrUploads {
			upload = util.CleanUrlPath(upload)
			param.Uploads[upload] = true
		}
	}

	// shows
	shows, err := getWildcardRegexp(result.GetValues("shows"))
	serverError.CheckFatal(err)
	param.Shows = shows

	showDirs, err := getWildcardRegexp(result.GetValues("showdirs"))
	serverError.CheckFatal(err)
	param.ShowDirs = showDirs

	showFiles, err := getWildcardRegexp(result.GetValues("showfiles"))
	serverError.CheckFatal(err)
	param.ShowFiles = showFiles

	// hides
	hides, err := getWildcardRegexp(result.GetValues("hides"))
	serverError.CheckFatal(err)
	param.Hides = hides

	hideDirs, err := getWildcardRegexp(result.GetValues("hidedirs"))
	serverError.CheckFatal(err)
	param.HideDirs = hideDirs

	hideFiles, err := getWildcardRegexp(result.GetValues("hidefiles"))
	serverError.CheckFatal(err)
	param.HideFiles = hideFiles
}

func Parse() *Param {
	paramCopied := param
	return &paramCopied
}
