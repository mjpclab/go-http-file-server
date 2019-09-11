package param

import (
	"regexp"
)

type Param struct {
	Root          string
	Aliases       map[string]string
	GlobalUpload  bool
	Uploads       []string
	GlobalArchive bool
	Archives      []string
	Key           string
	Cert          string
	Listen        string
	Template      string
	Shows         *regexp.Regexp
	ShowDirs      *regexp.Regexp
	ShowFiles     *regexp.Regexp
	Hides         *regexp.Regexp
	HideDirs      *regexp.Regexp
	HideFiles     *regexp.Regexp
	AccessLog     string
	ErrorLog      string
}
