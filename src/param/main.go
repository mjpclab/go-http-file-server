package param

import (
	"crypto/tls"
	"mjpclab.dev/ghfs/src/middleware"
	"mjpclab.dev/ghfs/src/serverError"
	"mjpclab.dev/ghfs/src/util"
	"os"
	"path/filepath"
)

type Param struct {
	Root      string
	EmptyRoot bool

	PrefixUrls   []string
	AutoDirSlash int

	DefaultSort string
	DirIndexes  []string
	// value: [url-path, fs-path]
	Aliases [][2]string

	GlobalRestrictAccess []string
	// value: [restrict-path, allow-hosts...]
	RestrictAccessUrls [][]string
	RestrictAccessDirs [][]string

	// value: [name, value]
	GlobalHeaders [][2]string
	// value: [path, (name, value)...]
	HeadersUrls [][]string
	HeadersDirs [][]string

	GlobalUpload bool
	UploadUrls   []string
	UploadDirs   []string

	GlobalMkdir bool
	MkdirUrls   []string
	MkdirDirs   []string

	GlobalDelete bool
	DeleteUrls   []string
	DeleteDirs   []string

	GlobalArchive bool
	ArchiveUrls   []string
	ArchiveDirs   []string

	GlobalCors bool
	CorsUrls   []string
	CorsDirs   []string

	GlobalAuth bool
	AuthUrls   []string
	AuthDirs   []string
	// value: [username, password]
	UsersPlain    [][2]string
	UsersBase64   [][2]string
	UsersMd5      [][2]string
	UsersSha1     [][2]string
	UsersSha256   [][2]string
	UsersSha512   [][2]string
	UserMatchCase bool

	Certificates []tls.Certificate
	Listens      []string
	ListensPlain []string
	ListensTLS   []string
	HostNames    []string
	Theme        string
	ThemeDir     string

	Hsts        bool
	HstsMaxAge  int
	ToHttps     bool
	ToHttpsPort string

	Shows     []string
	ShowDirs  []string
	ShowFiles []string
	Hides     []string
	HideDirs  []string
	HideFiles []string

	AccessLog string
	ErrorLog  string

	PreMiddlewares  []middleware.Middleware
	InMiddlewares   []middleware.Middleware
	PostMiddlewares []middleware.Middleware
}

type Params []*Param

func (param *Param) normalize() (errs []error) {
	var es []error
	var err error

	// root
	param.Root, err = filepath.Abs(param.Root)
	errs = serverError.AppendError(errs, err)

	// aliases
	param.Aliases, es = normalizePathMaps(param.Aliases)
	errs = append(errs, es...)

	// root & empty root && alias
	rootAliasIndex := -1
	for i := range param.Aliases {
		if param.Aliases[i][0] == "/" {
			rootAliasIndex = i
			break
		}
	}
	if rootAliasIndex >= 0 {
		param.EmptyRoot = false
	} else if param.EmptyRoot {
		param.Root = os.DevNull
		param.Aliases = append(param.Aliases, [2]string{"/", os.DevNull})
	} else {
		param.Aliases = append(param.Aliases, [2]string{"/", param.Root})
	}

	// url prefixes
	param.PrefixUrls = NormalizeUrlPaths(param.PrefixUrls)

	// // force dir slash
	if param.AutoDirSlash != 0 {
		param.AutoDirSlash = NormalizeRedirectCode(param.AutoDirSlash)
	}

	// dir indexes
	param.DirIndexes = normalizeFilenames(param.DirIndexes)

	// global restrict access, nil to disable, non-nil to enable with allowed hosts
	if param.GlobalRestrictAccess != nil {
		param.GlobalRestrictAccess = util.ExtractHostsFromUrls(param.GlobalRestrictAccess)
	}

	// restrict access
	param.RestrictAccessUrls, es = normalizeAllPathValues(param.RestrictAccessUrls, true, util.NormalizeUrlPath, util.ExtractHostsFromUrls)
	if len(es) == 0 {
		dedupAllPathValues(param.RestrictAccessUrls)
	} else {
		errs = append(errs, es...)
	}

	param.RestrictAccessDirs, es = normalizeAllPathValues(param.RestrictAccessDirs, true, filepath.Abs, util.ExtractHostsFromUrls)
	if len(es) == 0 {
		dedupAllPathValues(param.RestrictAccessDirs)
	} else {
		errs = append(errs, es...)
	}

	// headers
	param.HeadersUrls, es = normalizeAllPathValues(param.HeadersUrls, false, util.NormalizeUrlPath, normalizeHeaders)
	errs = append(errs, es...)

	param.HeadersDirs, es = normalizeAllPathValues(param.HeadersDirs, false, filepath.Abs, normalizeHeaders)
	errs = append(errs, es...)

	// upload/mkdir/delete/archive/cors/auth urls/dirs
	param.UploadUrls = NormalizeUrlPaths(param.UploadUrls)
	param.UploadDirs = NormalizeFsPaths(param.UploadDirs)
	param.MkdirUrls = NormalizeUrlPaths(param.MkdirUrls)
	param.MkdirDirs = NormalizeFsPaths(param.MkdirDirs)
	param.DeleteUrls = NormalizeUrlPaths(param.DeleteUrls)
	param.DeleteDirs = NormalizeFsPaths(param.DeleteDirs)
	param.ArchiveUrls = NormalizeUrlPaths(param.ArchiveUrls)
	param.ArchiveDirs = NormalizeFsPaths(param.ArchiveDirs)
	param.CorsUrls = NormalizeUrlPaths(param.CorsUrls)
	param.CorsDirs = NormalizeFsPaths(param.CorsDirs)
	param.AuthUrls = NormalizeUrlPaths(param.AuthUrls)
	param.AuthDirs = NormalizeFsPaths(param.AuthDirs)

	// hsts & https
	if param.Hsts {
		param.Hsts = validateHstsPort(param.ListensPlain, param.ListensTLS)
	}

	if param.ToHttps {
		param.ToHttpsPort, param.ToHttps = normalizeToHttpsPort(param.ToHttpsPort, param.ListensTLS)
	}

	return
}

func NewParams(paramList []Param) (params Params, errs []error) {
	params = make(Params, len(paramList))

	for i := range params {
		copiedParam := paramList[i]
		params[i] = &copiedParam
		errs = append(errs, params[i].normalize()...)
	}

	return
}
