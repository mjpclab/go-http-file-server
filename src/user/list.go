package user

import (
	"errors"
	"strings"
)

type List struct {
	users []user
}

func (list *List) Len() int {
	return len(list.users)
}

func (list *List) FindIndex(username string) int {
	for i := range list.users {
		if strings.EqualFold(list.users[i].getName(), username) {
			return i
		}
	}
	return -1
}

func (list *List) addUser(user user) error {
	username := user.getName()
	index := list.FindIndex(username)
	if index < 0 {
		list.users = append(list.users, user)
		return nil
	} else {
		return errors.New("duplicated username: " + username)
	}
}

func (list *List) AddPlain(username, password string) error {
	user := newPlainUser(username, password)
	err := list.addUser(user)
	return err
}

func (list *List) AddBase64(username, password string) error {
	user := newBase64User(username, password)
	err := list.addUser(user)
	return err
}

func (list *List) AddMd5(username, password string) error {
	user, err := newMd5User(username, password)
	if err != nil {
		return err
	}

	err = list.addUser(user)
	return err
}

func (list *List) AddSha1(username, password string) error {
	user, err := newSha1User(username, password)
	if err != nil {
		return err
	}

	err = list.addUser(user)
	return err
}

func (list *List) AddSha256(username, password string) error {
	user, err := newSha256User(username, password)
	if err != nil {
		return err
	}

	err = list.addUser(user)
	return err
}

func (list *List) AddSha512(username, password string) error {
	user, err := newSha512User(username, password)
	if err != nil {
		return err
	}

	err = list.addUser(user)
	return err
}

func (list *List) Auth(username, password string) (int, string, bool) {
	index := list.FindIndex(username)
	if index < 0 {
		return index, "", false
	}

	u := list.users[index]
	return index, u.getName(), u.auth(password)
}

func NewList() *List {
	return &List{}
}
