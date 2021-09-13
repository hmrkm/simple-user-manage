package adapter

import (
	"github.com/hmrkm/simple-user-manage/usecase"
)

type UsersAdapter interface {
	GetList(GetV1UsersParams) (ResponseGetUsers, error)
	Post(RequestPostUsers) error
	Get(string) (ResponseGetUser, error)
	Put(RequestPutUsers) error
	Delete(string) error
}

type usersAdapter struct {
	userService usecase.UserService
}

func NewUsersAdapter(us usecase.UserService) UsersAdapter {
	return usersAdapter{
		userService: us,
	}
}

func (u usersAdapter) GetList(req GetV1UsersParams) (ResponseGetUsers, error) {
	return ResponseGetUsers{}, nil
}

func (u usersAdapter) Post(req RequestPostUsers) error {
	return nil
}

func (u usersAdapter) Get(id string) (ResponseGetUser, error) {
	return ResponseGetUser{}, nil
}

func (u usersAdapter) Put(req RequestPutUsers) error {
	return nil
}

func (u usersAdapter) Delete(id string) error {
	return nil
}
