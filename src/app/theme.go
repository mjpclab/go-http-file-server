package app

import (
	"mjpclab.dev/ghfs/src/tpl/theme"
	"path/filepath"
)

func loadTheme(themePath string, themePool map[string]theme.Theme) (theme.Theme, []error) {
	themeKey, err := filepath.Abs(themePath)
	if err != nil {
		return nil, []error{err}
	}

	themeInst, themeExists := themePool[themeKey]
	if themeExists {
		return themeInst, nil
	}

	themeInst, err = theme.LoadMemTheme(themeKey)
	if err != nil {
		return nil, []error{err}
	}

	themePool[themeKey] = themeInst
	return themeInst, nil
}
