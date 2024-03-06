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
	Aliases     [][2]string // [][url-path, fs-path]

	// [][username, password]
	UsersPlain  [][2]string
	UsersBase64 [][2]string
	UsersMd5    [][2]string
	UsersSha1   [][2]string
	UsersSha256 [][2]string
	UsersSha512 [][2]string

	GlobalAuth    bool
	AuthUrls      []string
	AuthUrlsUsers [][]string // [][path, user...]
	AuthDirs      []string
	AuthDirsUsers [][]string // [][path, user...]

	GlobalRestrictAccess []string
	// [][restrict-path, allow-hosts...]
	RestrictAccessUrls [][]string
	RestrictAccessDirs [][]string

	GlobalHeaders [][2]string // [][name, value]
	// [][path, (name, value)...]
	HeadersUrls [][]string
	HeadersDirs [][]string

	GlobalUpload    bool
	UploadUrls      []string
	UploadUrlsUsers [][]string // [][path, user...]
	UploadDirs      []string
	UploadDirsUsers [][]string // [][path, user...]

	GlobalMkdir    bool
	MkdirUrls      []string
	MkdirUrlsUsers [][]string // [][path, user...]
	MkdirDirs      []string
	MkdirDirsUsers [][]string // [][path, user...]

	GlobalDelete    bool
	DeleteUrls      []string
	DeleteUrlsUsers [][]string // [][path, user...]
	DeleteDirs      []string
	DeleteDirsUsers [][]string // [][path, user...]

	GlobalArchive bool
	ArchiveUrls   []string
	ArchiveDirs   []string

	GlobalCors bool
	CorsUrls   []string
	CorsDirs   []string

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
	PostMiddlewares []middleware.Middleware
}

type Params []*Param

func (param *Param) normalize() (errs []error) {
	var es []error
	var err error

	// listens
	param.Listens = util.InPlaceDedup(param.Listens)
	param.ListensPlain = util.InPlaceDedup(param.ListensPlain)
	param.ListensTLS = util.InPlaceDedup(param.ListensTLS)

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

	// auth/upload/mkdir/delete/archive/cors urls/dirs
	param.AuthUrls = NormalizeUrlPaths(param.AuthUrls)
	param.AuthDirs = NormalizeFsPaths(param.AuthDirs)
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

	// auth/upload/mkdir/delete urls/dirs users
	param.AuthUrlsUsers, es = normalizeAllPathValues(param.AuthUrlsUsers, true, util.NormalizeUrlPath, nil)
	if len(es) == 0 {
		dedupAllPathValues(param.AuthUrlsUsers)
	} else {
		errs = append(errs, es...)
	}
	param.AuthDirsUsers, es = normalizeAllPathValues(param.AuthDirsUsers, true, filepath.Abs, nil)
	if len(es) == 0 {
		dedupAllPathValues(param.AuthDirsUsers)
	} else {
		errs = append(errs, es...)
	}

	param.UploadUrlsUsers, es = normalizeAllPathValues(param.UploadUrlsUsers, false, util.NormalizeUrlPath, nil)
	if len(es) == 0 {
		dedupAllPathValues(param.UploadUrlsUsers)
	} else {
		errs = append(errs, es...)
	}
	param.UploadDirsUsers, es = normalizeAllPathValues(param.UploadDirsUsers, false, filepath.Abs, nil)
	if len(es) == 0 {
		dedupAllPathValues(param.UploadDirsUsers)
	} else {
		errs = append(errs, es...)
	}

	param.MkdirUrlsUsers, es = normalizeAllPathValues(param.MkdirUrlsUsers, false, util.NormalizeUrlPath, nil)
	if len(es) == 0 {
		dedupAllPathValues(param.MkdirUrlsUsers)
	} else {
		errs = append(errs, es...)
	}
	param.MkdirDirsUsers, es = normalizeAllPathValues(param.MkdirDirsUsers, false, filepath.Abs, nil)
	if len(es) == 0 {
		dedupAllPathValues(param.MkdirDirsUsers)
	} else {
		errs = append(errs, es...)
	}

	param.DeleteUrlsUsers, es = normalizeAllPathValues(param.DeleteUrlsUsers, false, util.NormalizeUrlPath, nil)
	if len(es) == 0 {
		dedupAllPathValues(param.DeleteUrlsUsers)
	} else {
		errs = append(errs, es...)
	}
	param.DeleteDirsUsers, es = normalizeAllPathValues(param.DeleteDirsUsers, false, filepath.Abs, nil)
	if len(es) == 0 {
		dedupAllPathValues(param.DeleteDirsUsers)
	} else {
		errs = append(errs, es...)
	}

	// global restrict access, nil to disable, non-nil to enable with allowed hosts
	if param.GlobalRestrictAccess != nil {
		param.GlobalRestrictAccess = util.ExtractHostsFromUrls(param.GlobalRestrictAccess)
		param.GlobalRestrictAccess = util.InPlaceDedup(param.GlobalRestrictAccess)
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
