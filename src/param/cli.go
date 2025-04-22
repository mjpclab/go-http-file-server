package param

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"mjpclab.dev/ghfs/src/goNixArgParser"
	"mjpclab.dev/ghfs/src/goVirtualHost"
	"mjpclab.dev/ghfs/src/serverError"
)

var cliCmd = NewCliCmd()

func NewCliCmd() *goNixArgParser.Command {
	cmd := goNixArgParser.NewSimpleCommand(os.Args[0], "Simple command line based HTTP file server to share local file system")
	options := cmd.Options()
	var opt goNixArgParser.Option

	// define option
	var err error
	err = options.AddFlagsValue("root", []string{"-r", "--root"}, "GHFS_ROOT", ".", "root directory of server")
	serverError.CheckFatal(err)

	err = options.AddFlags("emptyroot", []string{"-R", "--empty-root"}, "GHFS_EMPTY_ROOT", "use virtual empty root directory")
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("aliases", []string{"-a", "--alias"}, "", nil, "set alias path, <sep><url-path><sep><fs-path>, e.g. :/doc:/usr/share/doc")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("prefixurls", "--prefix", "", nil, "serve files under URL path instead of /")
	serverError.CheckFatal(err)

	err = options.AddFlagsValue("autodirslash", []string{"-/", "--auto-dir-slash"}, "GHFS_AUTO_DIR_SLASH", "", "auto redirect directory with \"/\" suffix, or file without suffix")
	serverError.CheckFatal(err)

	opt = goNixArgParser.NewFlagValueOption("defaultsort", "--default-sort", "GHFS_DEFAULT_SORT", "/n", "default sort for files and directories")
	opt.Description = "Available sort key:\n- `n` sort by name ascending\n- `N` sort by name descending\n- `e` sort by type(suffix) ascending\n- `E` sort by type(suffix) descending\n- `s` sort by size ascending\n- `S` sort by size descending\n- `t` sort by modify time ascending\n- `T` sort by modify time descending\n- `_` no sort\nDirectory sort:\n- `/<key>` directories before files\n- `<key>/` directories after files\n- `<key>` directories mixed with files\n"
	err = options.Add(opt)
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("dirindexes", []string{"-I", "--dir-index"}, "GHFS_DIR_INDEX", nil, "default index page for directory")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("users", "--user", "", nil, "user info: <username>:<password>")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("usersbase64", "--user-base64", "", nil, "user info: <username>:<base64-password>")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("usersmd5", "--user-md5", "", nil, "user info: <username>:<md5-password>")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("userssha1", "--user-sha1", "", nil, "user info: <username>:<sha1-password>")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("userssha256", "--user-sha256", "", nil, "user info: <username>:<sha256-password>")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("userssha512", "--user-sha512", "", nil, "user info: <username>:<sha512-password>")
	serverError.CheckFatal(err)

	err = options.AddFlag("globalauth", "--global-auth", "GHFS_GLOBAL_AUTH", "require Basic Auth for all directories")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("authurls", "--auth", "", nil, "url path that require Basic Auth")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("authurlsusers", "--auth-user", "", nil, "url path that require Basic Auth for specific users, <sep><url-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("authdirs", "--auth-dir", "", nil, "file system path that require Basic Auth")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("authdirsusers", "--auth-dir-user", "", nil, "file system path that require Basic Auth for specific users, <sep><fs-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("indexurls", "--index", "", []string{"/"}, "url path that allow directory index")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("indexurlsusers", "--index-user", "", nil, "url path that allow index files for specific users, <sep><url-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("indexdirs", "--index-dir", "", nil, "file system path that allow index files")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("indexdirsusers", "--index-dir-user", "", nil, "file system path that allow index files for specific users, <sep><fs-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlags("globalupload", []string{"-U", "--global-upload"}, "", "allow upload files for all url paths")
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("uploadurls", []string{"-u", "--upload"}, "", nil, "url path that allow upload files")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("uploadurlsusers", "--upload-user", "", nil, "url path that allow upload files for specific users, <sep><url-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("uploaddirs", []string{"-p", "--upload-dir"}, "", nil, "file system path that allow upload files")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("uploaddirsusers", "--upload-dir-user", "", nil, "file system path that allow upload files for specific users, <sep><fs-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlag("globalmkdir", "--global-mkdir", "", "allow mkdir files for all url paths")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("mkdirurls", "--mkdir", "", nil, "url path that allow mkdir files")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("mkdirurlsusers", "--mkdir-user", "", nil, "url path that allow mkdir files for specific users, <sep><url-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("mkdirdirs", "--mkdir-dir", "", nil, "file system path that allow mkdir files")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("mkdirdirsusers", "--mkdir-dir-user", "", nil, "file system path that allow mkdir files for specific users, <sep><fs-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlag("globaldelete", "--global-delete", "", "allow delete files for all url paths")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("deleteurls", "--delete", "", nil, "url path that allow delete files")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("deleteurlsusers", "--delete-user", "", nil, "url path that allow delete files for specific users, <sep><url-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("deletedirs", "--delete-dir", "", nil, "file system path that allow delete files")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("deletedirsusers", "--delete-dir-user", "", nil, "file system path that allow delete files for specific users, <sep><fs-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlags("globalarchive", []string{"-A", "--global-archive"}, "GHFS_GLOBAL_ARCHIVE", "enable download archive for all directories")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("archiveurls", "--archive", "", nil, "url path that enable download as archive for specific directories")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("archiveurlsusers", "--archive-user", "", nil, "url path that allow archive files for specific users, <sep><url-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("archivedirs", "--archive-dir", "", nil, "file system path that enable download as archive for specific directories")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("archivedirsusers", "--archive-dir-user", "", nil, "file system path that allow archive files for specific users, <sep><fs-path>[<sep><user>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValue("maxarchiveworkers", "--max-archive-workers", "", "-1", "maximum number of concurrent archive operations (-1 for unlimited)")
	serverError.CheckFatal(err)

	err = options.AddFlag("globalcors", "--global-cors", "GHFS_GLOBAL_CORS", "enable CORS headers for all directories")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("corsurls", "--cors", "", nil, "url path that enable CORS headers")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("corsdirs", "--cors-dir", "", nil, "file system path that enable CORS headers")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("globalrestrictaccess", "--global-restrict-access", "GHFS_GLOBAL_RESTRICT_ACCESS", []string{}, "restrict access to all url paths from current host, with optional extra allow list")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("restrictaccessurls", "--restrict-access", "", []string{}, "restrict access to specific url paths from current host, with optional extra allow list, <sep><url-path>[<sep><allowed-host>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("restrictaccessdirs", "--restrict-access-dir", "", []string{}, "restrict access to specific file system paths from current host, with optional extra allow list, <sep><fs-path>[<sep><allowed-host>...]")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("globalheaders", "--global-header", "GHFS_GLOBAL_HEADER", []string{}, "custom headers for all url paths, e.g. <name>:<value>")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("headersurls", "--header", "", []string{}, "url path for custom headers, <sep><url><sep><name><sep><value>")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("headersdirs", "--header-dir", "", []string{}, "file system path for custom headers, <sep><dir><sep><name><sep><value>")
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("certs", []string{"-c", "--cert"}, "GHFS_CERT", nil, "TLS certificate path")
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("keys", []string{"-k", "--key"}, "GHFS_KEY", nil, "TLS certificate key path")
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("listens", []string{"-l", "--listen"}, "GHFS_LISTEN", nil, "address and port to listen")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("listensplain", "--listen-plain", "GHFS_LISTEN_PLAIN", nil, "address and port to listen, force plain http protocol")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("listenstls", "--listen-tls", "GHFS_LISTEN_TLS", nil, "address and port to listen, force https protocol")
	serverError.CheckFatal(err)

	err = options.AddFlagValues("hostnames", "--hostname", "", nil, "hostname for the virtual host")
	serverError.CheckFatal(err)

	err = options.AddFlagValue("theme", "--theme", "GHFS_THEME", "", "external theme file")
	serverError.CheckFatal(err)

	err = options.AddFlagValue("themedir", "--theme-dir", "GHFS_THEME_DIR", "", "external theme directory")
	serverError.CheckFatal(err)

	err = options.AddFlagValue("hsts", "--hsts", "GHFS_HSTS", "", "enable HSTS(HTTP Strict Transport Security)")
	serverError.CheckFatal(err)

	err = options.AddFlagValue("tohttps", "--to-https", "GHFS_TO_HTTPS", "", "redirect http:// to https://, with optional target port")
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("shows", []string{"-S", "--show"}, "GHFS_SHOW", nil, "show directories or files match wildcard")
	serverError.CheckFatal(err)
	err = options.AddFlagsValues("showdirs", []string{"-SD", "--show-dir"}, "GHFS_SHOW_DIR", nil, "show directories match wildcard")
	serverError.CheckFatal(err)
	err = options.AddFlagsValues("showfiles", []string{"-SF", "--show-file"}, "GHFS_SHOW_FILE", nil, "show files match wildcard")
	serverError.CheckFatal(err)

	err = options.AddFlagsValues("hides", []string{"-H", "--hide"}, "GHFS_HIDE", nil, "hide directories or files match wildcard")
	serverError.CheckFatal(err)
	err = options.AddFlagsValues("hidedirs", []string{"-HD", "--hide-dir"}, "GHFS_HIDE_DIR", nil, "hide directories match wildcard")
	serverError.CheckFatal(err)
	err = options.AddFlagsValues("hidefiles", []string{"-HF", "--hide-file"}, "GHFS_HIDE_FILE", nil, "hide files match wildcard")
	serverError.CheckFatal(err)

	err = options.AddFlagsValue("accesslog", []string{"-L", "--access-log"}, "GHFS_ACCESS_LOG", "", "access log file, use \"-\" for stdout")
	serverError.CheckFatal(err)

	err = options.AddFlagsValue("errorlog", []string{"-E", "--error-log"}, "GHFS_ERROR_LOG", "-", "error log file, use \"-\" for stderr")
	serverError.CheckFatal(err)

	err = options.AddFlagValue("config", "--config", "GHFS_CONFIG", "", "external config file")
	serverError.CheckFatal(err)

	err = options.AddFlag("version", "--version", "", "print version")
	serverError.CheckFatal(err)

	err = options.AddFlags("help", []string{"-h", "--help"}, "", "print this help")
	serverError.CheckFatal(err)

	return cmd
}

func ArgsToCmdResults(cmd *goNixArgParser.Command, args []string) (results []*goNixArgParser.ParseResult, printVersion, printHelp bool, errs []error) {
	// parse option
	results = cmd.ParseGroups(args, nil)

	// pre-check
	for _, result := range results {
		// undefined flags
		undefs := result.GetUndefs()
		if len(undefs) > 0 {
			errs = append(errs,
				errors.New("unknown option: "+strings.Join(undefs, " ")),
			)
		}

		// version
		if result.HasFlagKey("version") {
			printVersion = true
		}

		// help
		if result.HasFlagKey("help") {
			printHelp = true
		}
	}
	if printVersion || printHelp || len(errs) > 0 {
		return
	}

	// append config and re-parse
	configs := make([]string, 0, len(results))
	groupSeps := cmd.Options().GroupSeps()[0]
	hasConfig := false
	var stdinConfigArgs []string
	for i := range results {
		configs = append(configs, groupSeps)

		// config file
		config, _ := results[i].GetString("config")
		if len(config) == 0 {
			continue
		}

		var configArgs []string
		if stdinConfigArgs != nil && config == "-" {
			configArgs = stdinConfigArgs
		} else {
			var err error
			configArgs, err = goNixArgParser.LoadConfigArgs(config)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if config == "-" {
				stdinConfigArgs = configArgs
			}
		}
		if len(configArgs) == 0 {
			continue
		}

		hasConfig = true
		configs = append(configs, configArgs...)
	}
	if len(errs) > 0 {
		return
	}

	if hasConfig {
		configs = configs[1:]
		results = cmd.ParseGroups(args, configs)
		for i := range results {
			undefs := results[i].GetUndefs()
			if len(undefs) > 0 {
				errs = append(errs,
					errors.New("unknown option from config: "+strings.Join(undefs, " ")),
				)
			}
		}
		if len(errs) > 0 {
			return
		}
	}

	return
}

func CmdResultsToParams(results []*goNixArgParser.ParseResult) (params Params, errs []error) {
	var es []error

	// init param data
	params = make(Params, 0, len(results))
	for _, result := range results {
		param := &Param{}

		// regular option
		param.Root, _ = result.GetString("root")
		param.EmptyRoot = result.HasKey("emptyroot")
		param.PrefixUrls, _ = result.GetStrings("prefixurls")
		param.DefaultSort, _ = result.GetString("defaultsort")
		param.HostNames, _ = result.GetStrings("hostnames")
		param.Theme, _ = result.GetString("theme")
		param.ThemeDir, _ = result.GetString("themedir")
		param.AccessLog, _ = result.GetString("accesslog")
		param.ErrorLog, _ = result.GetString("errorlog")

		// aliases
		strAlias, _ := result.GetStrings("aliases")
		param.Aliases = SplitAllKeyValue(strAlias)

		// force dir slash
		if result.HasKey("autodirslash") {
			redirectCode, _ := result.GetInt("autodirslash")
			if redirectCode == 0 {
				redirectCode = http.StatusMovedPermanently
			}
			param.AutoDirSlash = redirectCode
		}

		// dir indexes
		param.DirIndexes, _ = result.GetStrings("dirindexes")

		// users
		arrUsersPlain, _ := result.GetStrings("users")
		param.UsersPlain = entriesToUsers(arrUsersPlain)
		arrUsersBase64, _ := result.GetStrings("usersbase64")
		param.UsersBase64 = entriesToUsers(arrUsersBase64)
		arrUsersMd5, _ := result.GetStrings("usersmd5")
		param.UsersMd5 = entriesToUsers(arrUsersMd5)
		arrUsersSha1, _ := result.GetStrings("userssha1")
		param.UsersSha1 = entriesToUsers(arrUsersSha1)
		arrUsersSha256, _ := result.GetStrings("userssha256")
		param.UsersSha256 = entriesToUsers(arrUsersSha256)
		arrUsersSha512, _ := result.GetStrings("userssha512")
		param.UsersSha512 = entriesToUsers(arrUsersSha512)

		// auth/index/upload/mkdir/delete/archive/cors urls/dirs
		param.GlobalAuth = result.HasKey("globalauth")
		param.AuthUrls, _ = result.GetStrings("authurls")
		param.AuthDirs, _ = result.GetStrings("authdirs")

		param.IndexUrls, _ = result.GetStrings("indexurls")
		param.IndexDirs, _ = result.GetStrings("indexdirs")

		param.GlobalUpload = result.HasKey("globalupload")
		param.UploadUrls, _ = result.GetStrings("uploadurls")
		param.UploadDirs, _ = result.GetStrings("uploaddirs")

		param.GlobalMkdir = result.HasKey("globalmkdir")
		param.MkdirUrls, _ = result.GetStrings("mkdirurls")
		param.MkdirDirs, _ = result.GetStrings("mkdirdirs")

		param.GlobalDelete = result.HasKey("globaldelete")
		param.DeleteUrls, _ = result.GetStrings("deleteurls")
		param.DeleteDirs, _ = result.GetStrings("deletedirs")

		param.GlobalArchive = result.HasKey("globalarchive")
		param.ArchiveUrls, _ = result.GetStrings("archiveurls")
		param.ArchiveDirs, _ = result.GetStrings("archivedirs")

		param.GlobalCors = result.HasKey("globalcors")
		param.CorsUrls, _ = result.GetStrings("corsurls")
		param.CorsDirs, _ = result.GetStrings("corsdirs")

		// auth/upload/mkdir/delete/archive urls/dirs urls users
		authUrlsUsers, _ := result.GetStrings("authurlsusers")
		param.AuthUrlsUsers = SplitAllKeyValues(authUrlsUsers)

		authDirsUsers, _ := result.GetStrings("authdirsusers")
		param.AuthDirsUsers = SplitAllKeyValues(authDirsUsers)

		indexUrlsUsers, _ := result.GetStrings("indexurlsusers")
		param.IndexUrlsUsers = SplitAllKeyValues(indexUrlsUsers)

		indexDirsUsers, _ := result.GetStrings("indexdirsusers")
		param.IndexDirsUsers = SplitAllKeyValues(indexDirsUsers)

		uploadUrlsUsers, _ := result.GetStrings("uploadurlsusers")
		param.UploadUrlsUsers = SplitAllKeyValues(uploadUrlsUsers)

		uploadDirsUsers, _ := result.GetStrings("uploaddirsusers")
		param.UploadDirsUsers = SplitAllKeyValues(uploadDirsUsers)

		mkdirUrlsUsers, _ := result.GetStrings("mkdirurlsusers")
		param.MkdirUrlsUsers = SplitAllKeyValues(mkdirUrlsUsers)

		mkdirDirsUsers, _ := result.GetStrings("mkdirdirsusers")
		param.MkdirDirsUsers = SplitAllKeyValues(mkdirDirsUsers)

		deleteUrlsUsers, _ := result.GetStrings("deleteurlsusers")
		param.DeleteUrlsUsers = SplitAllKeyValues(deleteUrlsUsers)

		deleteDirsUsers, _ := result.GetStrings("deletedirsusers")
		param.DeleteDirsUsers = SplitAllKeyValues(deleteDirsUsers)

		archiveUrlsUsers, _ := result.GetStrings("archiveurlsusers")
		param.ArchiveUrlsUsers = SplitAllKeyValues(archiveUrlsUsers)

		archiveDirsUsers, _ := result.GetStrings("archivedirsusers")
		param.ArchiveDirsUsers = SplitAllKeyValues(archiveDirsUsers)

		param.ArchiveMaxWorkers, _ = result.GetInt("maxarchiveworkers")

		// global restrict access
		if result.HasKey("globalrestrictaccess") {
			param.GlobalRestrictAccess, _ = result.GetStrings("globalrestrictaccess")
		}

		// restrict access urls
		restrictAccessUrls, _ := result.GetStrings("restrictaccessurls")
		param.RestrictAccessUrls = SplitAllKeyValues(restrictAccessUrls)

		// restrict access dirs
		restrictAccessDirs, _ := result.GetStrings("restrictaccessdirs")
		param.RestrictAccessDirs = SplitAllKeyValues(restrictAccessDirs)

		// global headers
		globalHeaders, _ := result.GetStrings("globalheaders")
		param.GlobalHeaders = EntriesToKVs(globalHeaders)

		// headers urls
		headersUrls, _ := result.GetStrings("headersurls")
		param.HeadersUrls = SplitAllKeyValues(headersUrls)

		// headers dirs
		headersDirs, _ := result.GetStrings("headersdirs")
		param.HeadersDirs = SplitAllKeyValues(headersDirs)

		// certificate
		certFiles, _ := result.GetStrings("certs")
		keyFiles, _ := result.GetStrings("keys")
		param.CertKeyPaths, _ = goVirtualHost.CertsKeysToPairs(certFiles, keyFiles)

		// listen
		listens, _ := result.GetStrings("listens")
		param.Listens = append(param.Listens, listens...)

		listenRests := result.GetRests()
		param.Listens = append(param.Listens, listenRests...)

		param.ListensPlain, _ = result.GetStrings("listensplain")

		param.ListensTLS, _ = result.GetStrings("listenstls")

		// hsts & https
		param.Hsts = result.HasKey("hsts")
		if param.Hsts {
			if result.HasValue("hsts") {
				param.HstsMaxAge, _ = result.GetInt("hsts")
			} else {
				param.HstsMaxAge = 31536000
			}
		}

		param.ToHttps = result.HasKey("tohttps")
		param.ToHttpsPort, _ = result.GetString("tohttps")

		// shows/hides
		param.Shows, _ = result.GetStrings("shows")
		param.ShowDirs, _ = result.GetStrings("showdirs")
		param.ShowFiles, _ = result.GetStrings("showfiles")
		param.Hides, _ = result.GetStrings("hides")
		param.HideDirs, _ = result.GetStrings("hidedirs")
		param.HideFiles, _ = result.GetStrings("hidefiles")

		es = param.Normalize()
		errs = append(errs, es...)

		params = append(params, param)
	}

	return
}

func ParseFromCli() (params Params, printVersion, printHelp bool, errs []error) {
	var cmdResults []*goNixArgParser.ParseResult

	cmdResults, printVersion, printHelp, errs = ArgsToCmdResults(cliCmd, os.Args)
	if printVersion || printHelp || len(errs) > 0 {
		return
	}

	params, errs = CmdResultsToParams(cmdResults)
	return
}

func PrintHelp() {
	cliCmd.OutputHelp(os.Stdout)
}
