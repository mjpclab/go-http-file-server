package param

import "strings"

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
