package param

import (
	argParser "../goNixArgParser"
	"../serverErrHandler"
	"../util"
	"os"
	"regexp"
	"strings"
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

	err = argParser.AddFlags("globalarchive", []string{"-A", "--global-archive"}, "GHFS_GLOBAL_ARCHIVE", "enable download archive for all directories")
	serverErrHandler.CheckFatal(err)

	err = argParser.AddFlagValues("archives", "--archive", "", nil, "enable download archive for specific directories")
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
	param.Root, _ = result.GetString("root")
	param.GlobalUpload = result.HasKey("globalupload")
	param.GlobalArchive = result.HasKey("globalarchive")
	param.Key, _ = result.GetString("key")
	param.Cert, _ = result.GetString("cert")
	if rests := result.GetRests(); len(rests) > 0 {
		param.Listen = rests[len(rests)-1]
	} else if listen, foundListen := result.GetString("listen"); foundListen {
		param.Listen = listen
	}
	param.Template, _ = result.GetString("template")
	param.AccessLog, _ = result.GetString("accesslog")
	param.ErrorLog, _ = result.GetString("errorlog")

	// normalize aliases
	arrAlias, _ := result.GetStrings("aliases")
	param.Aliases = normalizePathMaps(arrAlias)

	// normalize uploads
	arrUploads, _ := result.GetStrings("uploads")
	param.Uploads = normalizeUrlPaths(arrUploads)

	// normalize archives
	arrArchives, _ := result.GetStrings("archives")
	param.Archives = normalizeUrlPaths(arrArchives)

	// shows
	shows, err := getWildcardRegexp(result.GetStrings("shows"))
	serverErrHandler.CheckFatal(err)
	param.Shows = shows

	showDirs, err := getWildcardRegexp(result.GetStrings("showdirs"))
	serverErrHandler.CheckFatal(err)
	param.ShowDirs = showDirs

	showFiles, err := getWildcardRegexp(result.GetStrings("showfiles"))
	serverErrHandler.CheckFatal(err)
	param.ShowFiles = showFiles

	// hides
	hides, err := getWildcardRegexp(result.GetStrings("hides"))
	serverErrHandler.CheckFatal(err)
	param.Hides = hides

	hideDirs, err := getWildcardRegexp(result.GetStrings("hidedirs"))
	serverErrHandler.CheckFatal(err)
	param.HideDirs = hideDirs

	hideFiles, err := getWildcardRegexp(result.GetStrings("hidefiles"))
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
