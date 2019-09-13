package param

import (
	"../goNixArgParser"
	"../serverErrHandler"
	"io/ioutil"
	"os"
	"strings"
)

var cliParam *Param
var cliCmd *goNixArgParser.Command

func init() {
	cliCmd = goNixArgParser.NewSimpleCommand(os.Args[0], "Simple command line based HTTP file server to share local file system")
	optionSet := cliCmd.OptionSet

	// define option
	var err error
	err = optionSet.AddFlagsValue("root", []string{"-r", "--root"}, "GHFS_ROOT", ".", "root directory of server")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValues("aliases", []string{"-a", "--alias"}, "", nil, "set alias path, <sep><url><sep><path>, e.g. :/doc:/usr/share/doc")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlags("globalupload", []string{"-U", "--global-upload"}, "", "allow upload files for all url paths")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValues("uploads", []string{"-u", "--upload"}, "", nil, "url path that allow upload files")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlags("globalarchive", []string{"-A", "--global-archive"}, "GHFS_GLOBAL_ARCHIVE", "enable download archive for all directories")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagValues("archives", "--archive", "", nil, "enable download archive for specific directories")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValue("key", []string{"-k", "--key"}, "GHFS_KEY", "", "TLS certificate key path")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValue("cert", []string{"-c", "--cert"}, "GHFS_CERT", "", "TLS certificate path")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValue("listen", []string{"-l", "--listen"}, "GHFS_LISTEN", "", "address and port to listen")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValue("template", []string{"-t", "--template"}, "GHFS_TEMPLATE", "", "custom template file for page")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValues("shows", []string{"-S", "--show"}, "GHFS_SHOW", nil, "show directories or files match wildcard")
	serverErrHandler.CheckFatal(err)
	err = optionSet.AddFlagsValues("showdirs", []string{"-SD", "--show-dir"}, "GHFS_SHOW_DIR", nil, "show directories match wildcard")
	serverErrHandler.CheckFatal(err)
	err = optionSet.AddFlagsValues("showfiles", []string{"-SF", "--show-file"}, "GHFS_SHOW_FILE", nil, "show files match wildcard")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValues("hides", []string{"-H", "--hide"}, "GHFS_HIDE", nil, "hide directories or files match wildcard")
	serverErrHandler.CheckFatal(err)
	err = optionSet.AddFlagsValues("hidedirs", []string{"-HD", "--hide-dir"}, "GHFS_HIDE_DIR", nil, "hide directories match wildcard")
	serverErrHandler.CheckFatal(err)
	err = optionSet.AddFlagsValues("hidefiles", []string{"-HF", "--hide-file"}, "GHFS_HIDE_FILE", nil, "hide files match wildcard")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValue("accesslog", []string{"-L", "--access-log"}, "GHFS_ACCESS_LOG", "", "access log file, use \"-\" for stdout")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagsValue("errorlog", []string{"-E", "--error-log"}, "GHFS_ERROR_LOG", "-", "error log file, use \"-\" for stderr")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlagValue("config", "--config", "", "", "print this help")
	serverErrHandler.CheckFatal(err)

	err = optionSet.AddFlags("help", []string{"-h", "--help"}, "", "print this help")
	serverErrHandler.CheckFatal(err)
}

func doParseCli() *Param {
	param := &Param{}

	// parse option
	result := cliCmd.Parse(os.Args, nil)

	// help
	if result.HasFlagKey("help") {
		cliCmd.PrintHelp()
		os.Exit(0)
	}

	// config file
	if config, _ := result.GetString("config"); len(config) > 0 {
		configStr, err := ioutil.ReadFile(config)
		if !serverErrHandler.CheckError(err) && len(configStr) > 0 {
			configs := strings.Fields(string(configStr))
			if len(configs) > 0 {
				result = cliCmd.Parse(os.Args, configs)
			}
		}
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
