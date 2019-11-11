package user

import (
	"errors"
)

type Users map[string]user

func (users Users) checkExist(username string) error {
	if _, exist := users[username]; exist {
		return errors.New("username already exist")
	}
	return nil
}

func (users Users) AddPlain(username, password string) (err error) {
	err = users.checkExist(username);
	if err != nil {
		return
	}

	users[username] = newPlainUser(password)
	return
}

func (users Users) AddBase64(username, password string) (err error) {
	err = users.checkExist(username);
	if err != nil {
		return
	}

	users[username] = newBase64User(password)
	return
}

func (users Users) AddMd5(username, password string) (err error) {
	err = users.checkExist(username);
	if err != nil {
		return
	}

	user, err := newMd5User(password)
	if err != nil {
		return
	}

	users[username] = user
	return
}

func (users Users) AddSha1(username, password string) (err error) {
	err = users.checkExist(username);
	if err != nil {
		return
	}

	user, err := newSha1User(password)
	if err != nil {
		return
	}

	users[username] = user
	return
}

func (users Users) AddSha256(username, password string) (err error) {
	err = users.checkExist(username);
	if err != nil {
		return
	}

	user, err := newSha256User(password)
	if err != nil {
		return
	}

	users[username] = user
	return
}

func (users Users) AddSha512(username, password string) (err error) {
	err = users.checkExist(username);
	if err != nil {
		return
	}

	user, err := newSha512User(password)
	if err != nil {
		return
	}

	users[username] = user
	return
}

func (users Users) Auth(username, password string) bool {
	user, exist := users[username]
	return exist && user.auth(password)
}

func NewUsers() Users {
	return Users{}
}
