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

func getWildcardRegexp(wildcards []string, found bool) (*regexp.Regexp, error) {
	if !found || len(wildcards) == 0 {
		return nil, nil
	}

	normalizedWildcards := make([]string, 0, len(wildcards))
	for _, wildcard := range wildcards {
		if len(wildcard) == 0 {
			continue
		}
		normalizedWildcards = append(normalizedWildcards, util.WildcardToRegexp(wildcard))
	}

	if len(normalizedWildcards) == 0 {
		return nil, nil
	}

	exp := strings.Join(normalizedWildcards, "|")
	return regexp.Compile(exp)
}

func init() {
	argParser.CommandLine.Summary = "Simple command line based HTTP file server to share local file system"

	// define option
	var err error
	err = argParser.AddFlagsValue("root", []string{"-r", "--root"}, "GHFS_ROOT", ".", "root directory of server")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValues("aliases", []string{"-a", "--alias"}, "", nil, "set alias path, <sep><url><sep><path>, e.g. :/doc:/usr/share/doc")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValues("uploads", []string{"-u", "--upload"}, "", nil, "url path that allow upload files")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("key", []string{"-k", "--key"}, "GHFS_KEY", "", "TLS certificate key path")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("cert", []string{"-c", "--cert"}, "GHFS_CERT", "", "TLS certificate path")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("listen", []string{"-l", "--listen"}, "GHFS_LISTEN", "", "address and port to listen")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("template", []string{"-t", "--template"}, "GHFS_TEMPLATE", "", "custom template file for page")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValues("shows", []string{"-S", "--show"}, "GHFS_SHOW", nil, "show directories or files match wildcard")
	serverError.CheckFatal(err)
	err = argParser.AddFlagsValues("showdirs", []string{"-SD", "--show-dir"}, "GHFS_SHOW_DIR", nil, "show directories match wildcard")
	serverError.CheckFatal(err)
	err = argParser.AddFlagsValues("showfiles", []string{"-SF", "--show-file"}, "GHFS_SHOW_FILE", nil, "show files match wildcard")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValues("hides", []string{"-H", "--hide"}, "GHFS_HIDE", nil, "hide directories or files match wildcard")
	serverError.CheckFatal(err)
	err = argParser.AddFlagsValues("hidedirs", []string{"-HD", "--hide-dir"}, "GHFS_HIDE_DIR", nil, "hide directories match wildcard")
	serverError.CheckFatal(err)
	err = argParser.AddFlagsValues("hidefiles", []string{"-HF", "--hide-file"}, "GHFS_HIDE_FILE", nil, "hide files match wildcard")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("accesslog", []string{"-L", "--access-log"}, "GHFS_ACCESS_LOG", "", "access log file, use \"-\" for stdout")
	serverError.CheckFatal(err)

	err = argParser.AddFlagsValue("errorlog", []string{"-E", "--error-log"}, "GHFS_ERROR_LOG", "-", "error log file, use \"-\" for stderr")
	serverError.CheckFatal(err)

	err = argParser.AddFlags("help", []string{"-h", "--help"}, "print this help")
	serverError.CheckFatal(err)

	// parse option
	result := argParser.Parse()

	// help
	if result.HasFlagKey("help") {
		argParser.PrintHelp()
		os.Exit(0)
	}

	// normalize option
	param.Root, _ = result.GetValue("root")
	param.Key, _ = result.GetValue("key")
	param.Cert, _ = result.GetValue("cert")
	if rests := result.GetRests(); len(rests) > 0 {
		param.Listen = rests[len(rests)-1]
	} else if listen, foundListen := result.GetValue("listen"); foundListen {
		param.Listen = listen
	}
	param.Template, _ = result.GetValue("template")
	param.AccessLog, _ = result.GetValue("accesslog")
	param.ErrorLog, _ = result.GetValue("errorlog")

	// normalize aliases
	param.Aliases = map[string]string{}
	arrAlias, _ := result.GetValues("aliases")
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
	arrUploads, _ := result.GetValues("uploads")
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
