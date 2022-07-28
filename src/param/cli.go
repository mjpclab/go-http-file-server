package param

import (
	"../goNixArgParser"
	"../serverErrHandler"
	"../util"
	"../version"
	"fmt"
	"os"
	"strings"
)

var cliCmd *goNixArgParser.Command

func init() {
	cliCmd = goNixArgParser.NewSimpleCommand(os.Args[0], "Simple command line based HTTP file server to share local file system")
	options := cliCmd.Options()
	var opt goNixArgParser.Option

	// define option
	var err error
	err = options.AddFlagsValue("root", []string{"-r", "--root"}, "GHFS_ROOT", ".", "root directory of server")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlags("emptyroot", []string{"-R", "--empty-root"}, "GHFS_EMPTY_ROOT", "use virtual empty root directory")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("prefixurls", "--prefix", "", nil, "serve files under URL path instead of /")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValue("forcedirslash", []string{"-/", "--force-dir-slash"}, "GHFS_FORCE_DIR_SLASH", "", "auto redirect directory with \"/\" suffix")
	serverErrHandler.CheckFatal(err)

	opt = goNixArgParser.NewFlagValueOption("defaultsort", "--default-sort", "GHFS_DEFAULT_SORT", "/n", "default sort for files and directories")
	opt.Description = "Available sort key:\n- `n` sort by name ascending\n- `N` sort by name descending\n- `e` sort by type(suffix) ascending\n- `E` sort by type(suffix) descending\n- `s` sort by size ascending\n- `S` sort by size descending\n- `t` sort by modify time ascending\n- `T` sort by modify time descending\n- `_` no sort\nDirectory sort:\n- `/<key>` directories before files\n- `<key>/` directories after files\n- `<key>` directories mixed with files\n"
	err = options.Add(opt)
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("dirindexes", []string{"-I", "--dir-index"}, "GHFS_DIR_INDEX", nil, "default index page for directory")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("aliases", []string{"-a", "--alias"}, "", nil, "set alias path, <sep><url><sep><path>, e.g. :/doc:/usr/share/doc")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("globalrestrictaccess", "--global-restrict-access", "GHFS_GLOBAL_RESTRICT_ACCESS", []string{}, "restrict access to all url paths from current host, with optional extra allow list")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("restrictaccessurls", "--restrict-access", "", []string{}, "restrict access to specific url paths from current host, with optional extra allow list")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("restrictaccessdirs", "--restrict-access-dir", "", []string{}, "restrict access to specific file system paths from current host, with optional extra allow list")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("globalheaders", "--global-header", "GHFS_GLOBAL_HEADER", []string{}, "custom headers for all url paths, e.g. <name>:<value>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("headersurls", "--header", "", []string{}, "url path for custom headers, <sep><url><sep><name><sep><value>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("headersdirs", "--header-dir", "", []string{}, "file system path for custom headers, <sep><dir><sep><name><sep><value>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlags("globalupload", []string{"-U", "--global-upload"}, "", "allow upload files for all url paths")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("uploadurls", []string{"-u", "--upload"}, "", nil, "url path that allow upload files")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("uploaddirs", []string{"-p", "--upload-dir"}, "", nil, "file system path that allow upload files")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("globalmkdir", "--global-mkdir", "", "allow mkdir files for all url paths")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("mkdirurls", "--mkdir", "", nil, "url path that allow mkdir files")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("mkdirdirs", "--mkdir-dir", "", nil, "file system path that allow mkdir files")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("globaldelete", "--global-delete", "", "allow delete files for all url paths")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("deleteurls", "--delete", "", nil, "url path that allow delete files")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("deletedirs", "--delete-dir", "", nil, "file system path that allow delete files")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlags("globalarchive", []string{"-A", "--global-archive"}, "GHFS_GLOBAL_ARCHIVE", "enable download archive for all directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("archiveurls", "--archive", "", nil, "url path that enable download as archive for specific directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("archivedirs", "--archive-dir", "", nil, "file system path that enable download as archive for specific directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("globalcors", "--global-cors", "GHFS_GLOBAL_CORS", "enable CORS headers for all directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("corsurls", "--cors", "", nil, "url path that enable CORS headers")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("corsdirs", "--cors-dir", "", nil, "file system path that enable CORS headers")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("globalauth", "--global-auth", "GHFS_GLOBAL_AUTH", "require Basic Auth for all directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("authurls", "--auth", "", nil, "url path that require Basic Auth")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("authdirs", "--auth-dir", "", nil, "file system path that require Basic Auth")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("users", "--user", "", nil, "user info: <username>:<password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("usersbase64", "--user-base64", "", nil, "user info: <username>:<base64-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("usersmd5", "--user-md5", "", nil, "user info: <username>:<md5-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("userssha1", "--user-sha1", "", nil, "user info: <username>:<sha1-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("userssha256", "--user-sha256", "", nil, "user info: <username>:<sha256-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("userssha512", "--user-sha512", "", nil, "user info: <username>:<sha512-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("usermatchcase", "--user-match-case", "GHFS_USER_MATCH_CASE", "username should be case sensitive")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("certs", []string{"-c", "--cert"}, "GHFS_CERT", nil, "TLS certificate path")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("keys", []string{"-k", "--key"}, "GHFS_KEY", nil, "TLS certificate key path")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("listens", []string{"-l", "--listen"}, "GHFS_LISTEN", nil, "address and port to listen")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("listensplain", "--listen-plain", "GHFS_LISTEN_PLAIN", nil, "address and port to listen, force plain http protocol")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("listenstls", "--listen-tls", "GHFS_LISTEN_TLS", nil, "address and port to listen, force https protocol")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("hostnames", "--hostname", "", nil, "hostname for the virtual host")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValue("theme", "--theme", "GHFS_THEME", "", "external theme file")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValue("themedir", "--theme-dir", "GHFS_THEME_DIR", "", "external theme directory")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("globalhsts", "--hsts", "GHFS_HSTS", "enable HSTS(HTTP Strict Transport Security)")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValue("globalhttps", "--to-https", "GHFS_TO_HTTPS", "", "redirect http:// to https://, with optional target port")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("shows", []string{"-S", "--show"}, "GHFS_SHOW", nil, "show directories or files match wildcard")
	serverErrHandler.CheckFatal(err)
	err = options.AddFlagsValues("showdirs", []string{"-SD", "--show-dir"}, "GHFS_SHOW_DIR", nil, "show directories match wildcard")
	serverErrHandler.CheckFatal(err)
	err = options.AddFlagsValues("showfiles", []string{"-SF", "--show-file"}, "GHFS_SHOW_FILE", nil, "show files match wildcard")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("hides", []string{"-H", "--hide"}, "GHFS_HIDE", nil, "hide directories or files match wildcard")
	serverErrHandler.CheckFatal(err)
	err = options.AddFlagsValues("hidedirs", []string{"-HD", "--hide-dir"}, "GHFS_HIDE_DIR", nil, "hide directories match wildcard")
	serverErrHandler.CheckFatal(err)
	err = options.AddFlagsValues("hidefiles", []string{"-HF", "--hide-file"}, "GHFS_HIDE_FILE", nil, "hide files match wildcard")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValue("accesslog", []string{"-L", "--access-log"}, "GHFS_ACCESS_LOG", "", "access log file, use \"-\" for stdout")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValue("errorlog", []string{"-E", "--error-log"}, "GHFS_ERROR_LOG", "-", "error log file, use \"-\" for stderr")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValue("config", "--config", "GHFS_CONFIG", "", "external config file")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("version", "--version", "", "print version")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlags("help", []string{"-h", "--help"}, "", "print this help")
	serverErrHandler.CheckFatal(err)
}

func ParseCli() []*Param {
	args := os.Args

	// parse option
	results := cliCmd.ParseGroups(args, nil)

	// pre-check
	for _, result := range results {
		// undefined flags
		undefs := result.GetUndefs()
		if len(undefs) > 0 {
			fmt.Println("unknown option:", strings.Join(undefs, " "))
			os.Exit(1)
		}

		// version
		if result.HasFlagKey("version") {
			version.PrintVersion()
			os.Exit(0)
		}

		// help
		if result.HasFlagKey("help") {
			cliCmd.OutputHelp(os.Stdout)
			os.Exit(0)
		}
	}

	// append config and re-parse
	configs := []string{}
	groupSeps := cliCmd.Options().GroupSeps()[0]
	foundConfig := false
	for i := range results {
		configs = append(configs, groupSeps)

		// config file
		config, _ := results[i].GetString("config")
		if len(config) == 0 {
			continue
		}

		configStr, err := os.ReadFile(config)
		if serverErrHandler.CheckError(err) || len(configStr) == 0 {
			continue
		}

		configArgs := strings.Fields(string(configStr))
		if len(configArgs) == 0 {
			continue
		}

		foundConfig = true
		configs = append(configs, configArgs...)
	}

	if foundConfig {
		configs = configs[1:]
		results = cliCmd.ParseGroups(args, configs)
		for i := range results {
			undefs := results[i].GetUndefs()
			if len(undefs) > 0 {
				fmt.Println("unknown option from config:", strings.Join(undefs, " "))
				os.Exit(1)
			}
		}
	}

	// init param data
	params := make([]*Param, 0, len(results))
	var err error
	var errs []error
	for _, result := range results {
		param := &Param{}

		// normalize option
		param.Root, _ = result.GetString("root")
		param.EmptyRoot = result.HasKey("emptyroot")
		param.DefaultSort, _ = result.GetString("defaultsort")
		param.GlobalUpload = result.HasKey("globalupload")
		param.GlobalMkdir = result.HasKey("globalmkdir")
		param.GlobalDelete = result.HasKey("globaldelete")
		param.GlobalArchive = result.HasKey("globalarchive")
		param.GlobalCors = result.HasKey("globalcors")
		param.GlobalAuth = result.HasKey("globalauth")
		param.UserMatchCase = result.HasKey("usermatchcase")
		param.HostNames, _ = result.GetStrings("hostnames")
		param.Theme, _ = result.GetString("theme")
		param.ThemeDir, _ = result.GetString("themedir")
		param.AccessLog, _ = result.GetString("accesslog")
		param.ErrorLog, _ = result.GetString("errorlog")

		// root
		root, _ := result.GetString("root")
		root, err = util.NormalizeFsPath(root)
		serverErrHandler.CheckFatal(err)
		param.Root = root

		// normalize url prefixes
		prefixurls, _ := result.GetStrings("prefixurls")
		param.PrefixUrls = normalizeUrlPaths(prefixurls)

		// force dir slash
		if result.HasKey("forcedirslash") {
			forceDirSlash, _ := result.GetString("forcedirslash")
			param.ForceDirSlash = normalizeRedirectCode(forceDirSlash)
		}

		// dir indexes
		dirIndexes, _ := result.GetStrings("dirindexes")
		param.DirIndexes = normalizeFilenames(dirIndexes)

		// global restrict access
		if result.HasKey("globalrestrictaccess") {
			globalRestrictAccesses, _ := result.GetStrings("globalrestrictaccess")
			param.GlobalRestrictAccess = util.ExtractHostsFromUrls(globalRestrictAccesses)
		}

		// restrict access urls
		restrictAccessUrls, _ := result.GetStrings("restrictaccessurls")
		param.RestrictAccessUrls, errs = normalizePathRestrictAccesses(restrictAccessUrls, util.NormalizeUrlPath)
		serverErrHandler.CheckFatal(errs...)

		// restrict access dirs
		restrictAccessDirs, _ := result.GetStrings("restrictaccessdirs")
		param.RestrictAccessDirs, errs = normalizePathRestrictAccesses(restrictAccessDirs, util.NormalizeFsPath)
		serverErrHandler.CheckFatal(errs...)

		// global headers
		globalHeaders, _ := result.GetStrings("globalheaders")
		param.GlobalHeaders = entriesToHeaders(globalHeaders)

		// headers urls
		arrHeadersUrls, _ := result.GetStrings("headersurls")
		param.HeadersUrls, errs = normalizePathHeadersMap(arrHeadersUrls, util.NormalizeUrlPath)
		serverErrHandler.CheckFatal(errs...)

		// headers dirs
		arrHeadersDirs, _ := result.GetStrings("headersdirs")
		param.HeadersDirs, errs = normalizePathHeadersMap(arrHeadersDirs, util.NormalizeFsPath)
		serverErrHandler.CheckFatal(errs...)

		// certificate
		certFiles, _ := result.GetStrings("certs")
		keyFiles, _ := result.GetStrings("keys")
		certs, errs := LoadCertificates(certFiles, keyFiles)
		serverErrHandler.CheckFatal(errs...)
		param.Certificates = certs

		// normalize aliases
		arrAlias, _ := result.GetStrings("aliases")
		param.Aliases, errs = normalizePathMaps(arrAlias)
		serverErrHandler.CheckFatal(errs...)

		// normalize upload urls
		arrUploadUrls, _ := result.GetStrings("uploadurls")
		param.UploadUrls = normalizeUrlPaths(arrUploadUrls)

		// normalize upload dirs
		arrUploadDirs, _ := result.GetStrings("uploaddirs")
		param.UploadDirs = normalizeFsPaths(arrUploadDirs)

		// normalize mkdir urls
		arrMkdirUrls, _ := result.GetStrings("mkdirurls")
		param.MkdirUrls = normalizeUrlPaths(arrMkdirUrls)

		// normalize mkdir dirs
		arrMkdirDirs, _ := result.GetStrings("mkdirdirs")
		param.MkdirDirs = normalizeFsPaths(arrMkdirDirs)

		// normalize delete urls
		arrDeleteUrls, _ := result.GetStrings("deleteurls")
		param.DeleteUrls = normalizeUrlPaths(arrDeleteUrls)

		// normalize delete dirs
		arrDeleteDirs, _ := result.GetStrings("deletedirs")
		param.DeleteDirs = normalizeFsPaths(arrDeleteDirs)

		// normalize archive urls
		arrArchiveUrls, _ := result.GetStrings("archiveurls")
		param.ArchiveUrls = normalizeUrlPaths(arrArchiveUrls)

		// normalize archive dirs
		arrArchiveDirs, _ := result.GetStrings("archivedirs")
		param.ArchiveDirs = normalizeFsPaths(arrArchiveDirs)

		// normalize cors urls
		arrCorsUrls, _ := result.GetStrings("corsurls")
		param.CorsUrls = normalizeUrlPaths(arrCorsUrls)

		// normalize cors dirs
		arrCorsDirs, _ := result.GetStrings("corsdirs")
		param.CorsDirs = normalizeFsPaths(arrCorsDirs)

		// normalize auth urls
		arrAuthUrls, _ := result.GetStrings("authurls")
		param.AuthUrls = normalizeUrlPaths(arrAuthUrls)

		// normalize auth dirs
		arrAuthDirs, _ := result.GetStrings("authdirs")
		param.AuthDirs = normalizeFsPaths(arrAuthDirs)

		// normalize users
		arrUsersPlain, _ := result.GetStrings("users")
		param.UsersPlain = EntriesToUsers(arrUsersPlain)
		arrUsersBase64, _ := result.GetStrings("usersbase64")
		param.UsersBase64 = EntriesToUsers(arrUsersBase64)
		arrUsersMd5, _ := result.GetStrings("usersmd5")
		param.UsersMd5 = EntriesToUsers(arrUsersMd5)
		arrUsersSha1, _ := result.GetStrings("userssha1")
		param.UsersSha1 = EntriesToUsers(arrUsersSha1)
		arrUsersSha256, _ := result.GetStrings("userssha256")
		param.UsersSha256 = EntriesToUsers(arrUsersSha256)
		arrUsersSha512, _ := result.GetStrings("userssha512")
		param.UsersSha512 = EntriesToUsers(arrUsersSha512)

		dupUserNames := param.GetDupUserNames()
		if len(dupUserNames) > 0 {
			serverErrHandler.CheckFatal(fmt.Errorf("duplicated usernames: %q", dupUserNames))
		}

		// normalize listen
		listens, _ := result.GetStrings("listens")
		param.Listens = append(param.Listens, listens...)

		listenRests := result.GetRests()
		param.Listens = append(param.Listens, listenRests...)

		param.ListensPlain, _ = result.GetStrings("listensplain")

		param.ListensTLS, _ = result.GetStrings("listenstls")

		// hsts & https
		if len(param.ListensTLS) > 0 {
			if result.HasKey("globalhsts") {
				param.GlobalHsts = validateHstsPort(param.ListensPlain, param.ListensTLS)
			}

			if result.HasKey("globalhttps") {
				httpsPort, _ := result.GetString("globalhttps")
				param.HttpsPort, param.GlobalHttps = normalizeHttpsPort(httpsPort, param.ListensTLS)
			}
		}

		// shows
		shows, err := WildcardToRegexp(result.GetStrings("shows"))
		serverErrHandler.CheckFatal(err)
		param.Shows = shows

		showDirs, err := WildcardToRegexp(result.GetStrings("showdirs"))
		serverErrHandler.CheckFatal(err)
		param.ShowDirs = showDirs

		showFiles, err := WildcardToRegexp(result.GetStrings("showfiles"))
		serverErrHandler.CheckFatal(err)
		param.ShowFiles = showFiles

		// hides
		hides, err := WildcardToRegexp(result.GetStrings("hides"))
		serverErrHandler.CheckFatal(err)
		param.Hides = hides

		hideDirs, err := WildcardToRegexp(result.GetStrings("hidedirs"))
		serverErrHandler.CheckFatal(err)
		param.HideDirs = hideDirs

		hideFiles, err := WildcardToRegexp(result.GetStrings("hidefiles"))
		serverErrHandler.CheckFatal(err)
		param.HideFiles = hideFiles

		normalize(param)
		params = append(params, param)
	}

	return params
}
