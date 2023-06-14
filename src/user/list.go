package user

import (
	"errors"
	"mjpclab.dev/ghfs/src/util"
)

type List struct {
	users          []user
	namesEqualFunc util.StrEqualFunc
}

func (list *List) Len() int {
	return len(list.users)
}

func (list *List) findIndex(username string) int {
	for i := range list.users {
		if list.namesEqualFunc(list.users[i].getName(), username) {
			return i
		}
	}
	return -1
}

func (list *List) addUser(user user) error {
	username := user.getName()
	index := list.findIndex(username)
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

func (list *List) Auth(username, password string) bool {
	index := list.findIndex(username)
	if index < 0 {
		return false
	}

	return list.users[index].auth(password)
}

func NewList(nameCaseSensitive bool) *List {
	var namesEqualFunc util.StrEqualFunc
	if nameCaseSensitive {
		namesEqualFunc = util.IsStrEqualAccurate
	} else {
		namesEqualFunc = util.IsStrEqualNoCase
	}
	return &List{namesEqualFunc: namesEqualFunc}
}
