package adapter

import (
	"github.com/hmrkm/simple-user-manage/usecase"
)

type Users interface {
	List(RequestUsersList) (ResponseUsersList, error)
	Create(RequestUsersCreate) error
	Detail(RequestUsersDetail) (ResponseUsersDetail, error)
	Update(RequestUsersUpdate) error
	Delete(RequestUsersDelete) error
}

type users struct {
	users usecase.Users
}

func NewUsers(u usecase.Users) Users {
	return users{
		users: u,
	}
}

func (u users) List(req RequestUsersList) (ResponseUsersList, error) {
	list, now, last, err := u.users.List(req.Page, req.Limit)
	if err != nil {
		return ResponseUsersList{}, err
	}

	ul := []User{}
	for _, v := range list {
		ul = append(ul, User{
			Email: v.Email,
			Id:    v.Id,
		})
	}

	return ResponseUsersList{
		List: ul,
		Page: Page{
			Last: last,
			Now:  now,
		},
	}, nil
}

func (u users) Create(req RequestUsersCreate) error {
	return nil
}

func (u users) Detail(req RequestUsersDetail) (ResponseUsersDetail, error) {
	return ResponseUsersDetail{}, nil
}

func (u users) Update(req RequestUsersUpdate) error {
	return nil
}

func (u users) Delete(RequestUsersDelete) error {
	return nil
}
