package param

import (
	"../util"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

func splitMapping(input string) (k, v string, ok bool) {
	sep, sepLen := utf8.DecodeRuneInString(input)
	if sepLen == 0 {
		return
	}
	entry := input[sepLen:]
	if len(entry) == 0 {
		return
	}

	sepIndex := strings.IndexRune(entry, sep)
	if sepIndex <= 0 || sepIndex+sepLen == len(entry) {
		return
	}

	k = entry[:sepIndex]
	v = entry[sepIndex+sepLen:]
	return k, v, true
}

func normalizePathMaps(inputs []string) map[string]string {
	maps := map[string]string{}

	for _, input := range inputs {
		urlPath, fsPath, ok := splitMapping(input)
		if !ok {
			continue
		}

		urlPath = util.CleanUrlPath(urlPath)
		fsPath = path.Clean(fsPath)
		maps[urlPath] = fsPath
	}

	return maps
}

func normalizeProxyMaps(inputs []string) map[string]string {
	maps := map[string]string{}

	for _, input := range inputs {
		urlPath, target, ok := splitMapping(input)
		if !ok {
			continue
		}

		cleanUrlPath := util.CleanUrlPath(urlPath)
		if len(cleanUrlPath) > 1 && urlPath[len(urlPath)-1] == '/' {
			cleanUrlPath += "/"
		}
		maps[cleanUrlPath] = target
	}

	return maps
}

func normalizeUrlPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}
		outputs = append(outputs, util.CleanUrlPath(input))
	}

	return outputs
}

func normalizeFsPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}

		abs, err := filepath.Abs(input)
		if err != nil {
			continue
		}

		outputs = append(outputs, abs)
	}

	return outputs
}

func getUsers(userEntries []string) []*user {
	users := make([]*user, 0, len(userEntries))
	for _, userEntry := range userEntries {
		username := userEntry
		password := ""

		colonIndex := strings.IndexByte(userEntry, ':')
		if colonIndex >= 0 {
			username = userEntry[:colonIndex]
			password = userEntry[colonIndex+1:]
		}

		users = append(users, &user{username, password})
	}
	return users
}

func getDupUserNames(usersGroups ...[]*user) []string {
	userMap := map[string]bool{}
	dupUserMap := map[string]bool{}

	for _, users := range usersGroups {
		for _, user := range users {
			if userMap[user.Username] {
				dupUserMap[user.Username] = true
			}
			userMap[user.Username] = true
		}
	}

	dupUsers := make([]string, 0, len(dupUserMap))
	for username, _ := range dupUserMap {
		dupUsers = append(dupUsers, username)
	}
	return dupUsers
}

func getWildcardRegexp(wildcards []string, found bool) (*regexp.Regexp, error) {
	if !found || len(wildcards) == 0 {
		return nil, nil
	}

	normalizedWildcards := make([]string, 0, len(wildcards))
	for _, wildcard := range wildcards {
		if len(wildcard) == 0 {
			continue
		}
		normalizedWildcards = append(normalizedWildcards, util.WildcardToRegexp(wildcard))
	}

	if len(normalizedWildcards) == 0 {
		return nil, nil
	}

	exp := strings.Join(normalizedWildcards, "|")
	return regexp.Compile(exp)
}
