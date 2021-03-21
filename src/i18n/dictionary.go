package i18n

type Dictionary struct {
	Lang  string
	Trans *Translation
}

var Dictionaries = [...]Dictionary{
	{"en-us", &translationEnUs},
	{"en", &translationEnUs},
	{"zh-cn", &translationZhSimp},
	{"zh-tw", &translationZhTrad},
	{"zh-hk", &translationZhTrad},
	{"zh", &translationZhSimp},
}

var LanguageTags []string

func init() {
	count := len(Dictionaries)
	LanguageTags = make([]string, count)
	for i := 0; i < count; i++ {
		LanguageTags[i] = Dictionaries[i].Lang
	}
}
