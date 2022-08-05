package param

import (
	"../serverError"
	"../util"
	"crypto/tls"
	"os"
)

type user struct {
	Username string
	Password string
}

type Param struct {
	Root      string
	EmptyRoot bool

	PrefixUrls    []string
	ForceDirSlash int

	DefaultSort string
	DirIndexes  []string
	Aliases     map[string]string

	GlobalRestrictAccess []string
	RestrictAccessUrls   map[string][]string
	RestrictAccessDirs   map[string][]string

	// value: [name, value]
	GlobalHeaders [][2]string
	HeadersUrls   map[string][][2]string
	HeadersDirs   map[string][][2]string

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

	GlobalAuth    bool
	AuthUrls      []string
	AuthDirs      []string
	UsersPlain    []*user
	UsersBase64   []*user
	UsersMd5      []*user
	UsersSha1     []*user
	UsersSha256   []*user
	UsersSha512   []*user
	UserMatchCase bool

	Certificates []tls.Certificate
	Listens      []string
	ListensPlain []string
	ListensTLS   []string
	HostNames    []string
	Theme        string
	ThemeDir     string

	GlobalHsts  bool
	GlobalHttps bool
	HttpsPort   string

	Shows     []string
	ShowDirs  []string
	ShowFiles []string
	Hides     []string
	HideDirs  []string
	HideFiles []string

	AccessLog string
	ErrorLog  string
}

func (param *Param) normalize() (errs []error) {
	var err error

	// root
	param.Root, err = util.NormalizeFsPath(param.Root)
	errs = serverError.AppendError(errs, err)

	// root & empty root && alias
	_, hasRootAlias := param.Aliases["/"]
	if hasRootAlias {
		param.EmptyRoot = false
	} else if param.EmptyRoot {
		param.Root = os.DevNull
		param.Aliases["/"] = os.DevNull
	} else {
		param.Aliases["/"] = param.Root
	}

	// url prefixes
	param.PrefixUrls = normalizeUrlPaths(param.PrefixUrls)

	// // force dir slash
	if param.ForceDirSlash != 0 {
		param.ForceDirSlash = normalizeRedirectCode(param.ForceDirSlash)
	}

	// dir indexes
	param.DirIndexes = normalizeFilenames(param.DirIndexes)

	// global restrict access, nil to disable, non-nil to enable with allowed hosts
	if param.GlobalRestrictAccess != nil {
		param.GlobalRestrictAccess = util.ExtractHostsFromUrls(param.GlobalRestrictAccess)
	}

	// upload/mkdir/delete/archive/cors/auth urls/dirs
	param.UploadUrls = normalizeUrlPaths(param.UploadUrls)
	param.UploadDirs = normalizeFsPaths(param.UploadDirs)
	param.MkdirUrls = normalizeUrlPaths(param.MkdirUrls)
	param.MkdirDirs = normalizeFsPaths(param.MkdirDirs)
	param.DeleteUrls = normalizeUrlPaths(param.DeleteUrls)
	param.DeleteDirs = normalizeFsPaths(param.DeleteDirs)
	param.ArchiveUrls = normalizeUrlPaths(param.ArchiveUrls)
	param.ArchiveDirs = normalizeFsPaths(param.ArchiveDirs)
	param.CorsUrls = normalizeUrlPaths(param.CorsUrls)
	param.CorsDirs = normalizeFsPaths(param.CorsDirs)
	param.AuthUrls = normalizeUrlPaths(param.AuthUrls)
	param.AuthDirs = normalizeFsPaths(param.AuthDirs)

	return
}
