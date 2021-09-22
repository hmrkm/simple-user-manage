package usecase

import "github.com/hmrkm/simple-user-manage/domain"

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-user-manage/$GOPACKAGE -package=$GOPACKAGE -destination=users_mock.go
type Users interface {
	List(page int, limit int) (users []domain.User, now int, last int, err error)
	Create(email string, password string, passwordConf string) error
	Detail(id string) (domain.User, error)
	Update(id string, email string, password *string, passwordConf *string) error
	Delete(id string) error
}

type users struct {
	userService domain.UserService
	store       domain.Store
}

func NewUsers(us domain.UserService, s domain.Store) Users {
	return users{
		userService: us,
		store:       s,
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

func (u users) Create(email string, password string, passwordConf string) error {
	if !u.userService.VerifyPassword(password, passwordConf) {
		return domain.ErrInvalidValue
	}

	return u.userService.Create(
		domain.CreateULID(),
		email,
		domain.CreateHash(password),
	)
}

func (u users) Detail(id string) (domain.User, error) {
	return u.userService.Read(id)
}

func (u users) Update(id string, email string, password *string, passwordConf *string) error {
	hasPassword := password != nil && passwordConf != nil
	if hasPassword {
		if !u.userService.VerifyPassword(*password, *passwordConf) {
			return domain.ErrInvalidValue
		}
	}

	user, err := u.userService.Read(id)
	if err != nil {
		return err
	}

	if hasPassword {
		return user.UpdateWithPassword(u.store, domain.User{
			Email:          email,
			HashedPassword: domain.CreateHash(*password),
		})
	} else {
		return user.Update(u.store, domain.User{
			Email: email,
		})
	}

}

func (u users) Delete(id string) error {
	user, err := u.userService.Read(id)
	if err != nil {
		return err
	}

	return user.Delete(u.store)
}
