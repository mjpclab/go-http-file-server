package user

import (
	"errors"
)

type List struct {
	users []user
}

var errUserExists = errors.New("username already exist")

func (list *List) findIndex(username string) int {
	for i := range list.users {
		if list.users[i].getName() == username {
			return i
		}
	}
	return -1
}

func (list *List) AddPlain(username, password string) error {
	if list.findIndex(username) >= 0 {
		return errUserExists
	}

	user := newPlainUser(username, password)
	list.users = append(list.users, user)
	return nil
}

func (list *List) AddBase64(username, password string) error {
	if list.findIndex(username) >= 0 {
		return errUserExists
	}

	user := newBase64User(username, password)
	list.users = append(list.users, user)
	return nil
}

func (list *List) AddMd5(username, password string) error {
	if list.findIndex(username) >= 0 {
		return errUserExists
	}

	user, err := newMd5User(username, password)
	if err != nil {
		return err
	}

	list.users = append(list.users, user)
	return nil
}

func (list *List) AddSha1(username, password string) error {
	if list.findIndex(username) >= 0 {
		return errUserExists
	}

	user, err := newSha1User(username, password)
	if err != nil {
		return err
	}

	list.users = append(list.users, user)
	return nil
}

func (list *List) AddSha256(username, password string) error {
	if list.findIndex(username) >= 0 {
		return errUserExists
	}

	user, err := newSha256User(username, password)
	if err != nil {
		return err
	}

	list.users = append(list.users, user)
	return nil
}

func (list *List) AddSha512(username, password string) error {
	if list.findIndex(username) >= 0 {
		return errUserExists
	}

	user, err := newSha512User(username, password)
	if err != nil {
		return err
	}

	list.users = append(list.users, user)
	return nil
}

func (list *List) Auth(username, password string) bool {
	index := list.findIndex(username)
	if index < 0 {
		return false
	}

	return list.users[index].auth(password)
}

func NewList() *List {
	return &List{}
}
