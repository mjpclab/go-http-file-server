package user

import "errors"

type user struct {
	username string
	password string
}

func newUser(username, password string) *user {
	return &user{
		username,
		password,
	}
}

type Users struct {
	pool map[string]*user
}

func (users Users) Add(username, password string) error {
	if _, exist := users.pool[username]; exist {
		return errors.New("username already exist")
	}

	users.pool[username] = newUser(username, password)
	return nil
}

func (users Users) Auth(username, password string) bool {
	user, exist := users.pool[username]
	return exist && user.password == password
}

func NewUsers() Users {
	return Users{
		pool: map[string]*user{},
	}
}
