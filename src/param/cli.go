package param

import (
	argParser "../goNixArgParser"
	"../serverErrHandler"
	"../util"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"
)

var cliParam *Param

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
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValues("aliases", []string{"-a", "--alias"}, "", nil, "set alias path, <sep><url><sep><path>, e.g. :/doc:/usr/share/doc")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlags("globalupload", []string{"-U", "--global-upload"}, "", "allow upload files for all url paths")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValues("uploads", []string{"-u", "--upload"}, "", nil, "url path that allow upload files")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlags("archive", []string{"-A", "--archive"}, "GHFS_ARCHIVE", "enable download archive of current directory")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValue("key", []string{"-k", "--key"}, "GHFS_KEY", "", "TLS certificate key path")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValue("cert", []string{"-c", "--cert"}, "GHFS_CERT", "", "TLS certificate path")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValue("listen", []string{"-l", "--listen"}, "GHFS_LISTEN", "", "address and port to listen")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValue("template", []string{"-t", "--template"}, "GHFS_TEMPLATE", "", "custom template file for page")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValues("shows", []string{"-S", "--show"}, "GHFS_SHOW", nil, "show directories or files match wildcard")
	serverErrHandler.CheckFatal(err)
	err = argParser.AddFlagsValues("showdirs", []string{"-SD", "--show-dir"}, "GHFS_SHOW_DIR", nil, "show directories match wildcard")
	serverErrHandler.CheckFatal(err)
	err = argParser.AddFlagsValues("showfiles", []string{"-SF", "--show-file"}, "GHFS_SHOW_FILE", nil, "show files match wildcard")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValues("hides", []string{"-H", "--hide"}, "GHFS_HIDE", nil, "hide directories or files match wildcard")
	serverErrHandler.CheckFatal(err)
	err = argParser.AddFlagsValues("hidedirs", []string{"-HD", "--hide-dir"}, "GHFS_HIDE_DIR", nil, "hide directories match wildcard")
	serverErrHandler.CheckFatal(err)
	err = argParser.AddFlagsValues("hidefiles", []string{"-HF", "--hide-file"}, "GHFS_HIDE_FILE", nil, "hide files match wildcard")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValue("accesslog", []string{"-L", "--access-log"}, "GHFS_ACCESS_LOG", "", "access log file, use \"-\" for stdout")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagsValue("errorlog", []string{"-E", "--error-log"}, "GHFS_ERROR_LOG", "-", "error log file, use \"-\" for stderr")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlags("help", []string{"-h", "--help"}, "", "print this help")
	serverErrHandler.CheckFatal(err)
}

func doParseCli() *Param {
	param := &Param{}

	// parse option
	result := argParser.Parse()

	// help
	if result.HasFlagKey("help") {
		argParser.PrintHelp()
		os.Exit(0)
	}

	// normalize option
	param.Root, _ = result.GetValue("root")
	param.GlobalUpload = result.HasKey("globalupload")
	param.CanArchive = result.HasKey("archive")
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
			if len(alias) == 0 {
				continue
			}

			sepIndex := strings.Index(alias, string(sep))
			if sepIndex == -1 {
				continue
			}

			urlPath := alias[:sepIndex]
			if len(urlPath) == 0 {
				continue
			}
			urlPath = util.CleanUrlPath(urlPath)

			fsPath := alias[sepIndex+sepLen:]
			if len(fsPath) == 0 {
				continue
			}
			fsPath = path.Clean(fsPath)

			param.Aliases[urlPath] = fsPath
		}
	}

	// normalize uploads
	uploadArgValues, _ := result.GetValues("uploads")
	param.Uploads = make([]string, len(uploadArgValues))
	for i, upload := range uploadArgValues {
		param.Uploads[i] = util.CleanUrlPath(upload)
	}

	// shows
	shows, err := getWildcardRegexp(result.GetValues("shows"))
	serverErrHandler.CheckFatal(err)
	param.Shows = shows

	showDirs, err := getWildcardRegexp(result.GetValues("showdirs"))
	serverErrHandler.CheckFatal(err)
	param.ShowDirs = showDirs

	showFiles, err := getWildcardRegexp(result.GetValues("showfiles"))
	serverErrHandler.CheckFatal(err)
	param.ShowFiles = showFiles

	// hides
	hides, err := getWildcardRegexp(result.GetValues("hides"))
	serverErrHandler.CheckFatal(err)
	param.Hides = hides

	hideDirs, err := getWildcardRegexp(result.GetValues("hidedirs"))
	serverErrHandler.CheckFatal(err)
	param.HideDirs = hideDirs

	hideFiles, err := getWildcardRegexp(result.GetValues("hidefiles"))
	serverErrHandler.CheckFatal(err)
	param.HideFiles = hideFiles

	return param
}

func ParseCli() *Param {
	if cliParam == nil {
		cliParam = doParseCli()
	}

	paramCopied := *cliParam
	return &paramCopied
}
