package param

import (
	"../goVirtualHost"
	"../util"
	"crypto/tls"
	"regexp"
	"strings"
)

func LoadCertificates(certFiles, keyFiles []string) ([]tls.Certificate, []error) {
	return goVirtualHost.LoadCertificates(certFiles, keyFiles)
}

func EntriesToUsers(entries []string) []*user {
	users := make([]*user, 0, len(entries))
	for _, userEntry := range entries {
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

func entriesToHeaders(entries []string) [][2]string {
	headers := make([][2]string, 0, len(entries))
	for _, entry := range entries {
		colonIndex := strings.IndexByte(entry, ':')
		if colonIndex <= 0 || colonIndex == len(entry)-1 {
			continue
		}
		key := entry[:colonIndex]
		value := entry[colonIndex+1:]
		headers = append(headers, [2]string{key, value})
	}
	return headers
}

func WildcardToRegexp(wildcards []string, found bool) (*regexp.Regexp, error) {
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

func (param *Param) GetDupUserNames() []string {
	userNameEqual := util.IsStrEqualNoCase
	if param.UserMatchCase {
		userNameEqual = util.IsStrEqualAccurate
	}

	usersGroups := [][]*user{
		param.UsersPlain,
		param.UsersBase64,
		param.UsersMd5,
		param.UsersSha1,
		param.UsersSha256,
		param.UsersSha512,
	}

	checkedUsernames := make([]string, 0,
		len(param.UsersPlain)+
			len(param.UsersBase64)+
			len(param.UsersMd5)+
			len(param.UsersSha1)+
			len(param.UsersSha256)+
			len(param.UsersSha512),
	)
	var dupUserNames []string

	for _, users := range usersGroups {
	eachUser:
		for _, user := range users {
			for _, existingUsername := range checkedUsernames {
				if userNameEqual(user.Username, existingUsername) {
					if !util.Contains(dupUserNames, existingUsername) {
						dupUserNames = append(dupUserNames, existingUsername)
					}
					continue eachUser
				}
			}
			checkedUsernames = append(checkedUsernames, user.Username)
		}
	}

	return dupUserNames
}
