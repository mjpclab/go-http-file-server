package param

import (
	"../serverError"
	"../util"
	"crypto/tls"
	"os"
)

type Param struct {
	Root      string
	EmptyRoot bool

	PrefixUrls    []string
	ForceDirSlash int

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
	var es []error
	var err error

	// root
	param.Root, err = util.NormalizeFsPath(param.Root)
	errs = serverError.AppendError(errs, err)

	// alias
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

	// restrict access
	param.RestrictAccessUrls, es = normalizePathValues(param.RestrictAccessUrls, true, util.NormalizeUrlPath, util.ExtractHostsFromUrls)
	errs = append(errs, es...)

	param.RestrictAccessDirs, es = normalizePathValues(param.RestrictAccessDirs, true, util.NormalizeFsPath, util.ExtractHostsFromUrls)
	errs = append(errs, es...)

	// headers
	param.HeadersUrls, es = normalizePathValues(param.HeadersUrls, false, util.NormalizeUrlPath, normalizeHeaders)
	errs = append(errs, es...)

	param.HeadersDirs, es = normalizePathValues(param.HeadersDirs, false, util.NormalizeFsPath, normalizeHeaders)
	errs = append(errs, es...)

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
