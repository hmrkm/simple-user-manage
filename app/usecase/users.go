package usecase

import "github.com/hmrkm/simple-user-manage/domain"

type Users interface {
	List(page int, limit int) (users []domain.User, now int, last int, err error)
	Create()
	Detail()
	Update()
	Delete()
}

type users struct {
	userService domain.UserService
}

func NewUsers(us domain.UserService) Users {
	return users{
		userService: us,
	}
}

func (u users) List(page int, limit int) (users []domain.User, now int, last int, err error) {
	us, err := u.userService.ReadList(page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	c, err := u.userService.Count()
	if err != nil {
		return nil, 0, 0, err
	}

	now, last = domain.Pager(c, page, limit)

	return us, now, last, nil
}

func (u users) Create() {

}

func (u users) Detail() {

}

func (u users) Update() {

}

func (u users) Delete() {

}
