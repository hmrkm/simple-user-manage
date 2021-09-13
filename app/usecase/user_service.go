package usecase

import (
	"github.com/pkg/errors"
)

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-user-manage/$GOPACKAGE -package=$GOPACKAGE -destination=user_service_mock.go
type UserService interface {
	Create(id string, email string, hashedPassword string) error
	Read(id string) (User, error)
	ReadList(page int, limit int) ([]User, error)
	UpdateWithPassword(id string, email string, hashedPassword string) error
	Update(id string, email string) error
	Delete(id string) error
}

type userService struct {
	store Store
}

func NewUserService(s Store) UserService {
	return userService{
		store: s,
	}
}

func (us userService) Create(id string, email string, hashedPassword string) error {
	u := User{
		Id:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	return us.store.Create(&u)
}

func (us userService) Read(id string) (User, error) {
	u := User{}
	if err := us.store.First(&u, "id=?", id); err != nil {
		return User{}, errors.WithStack(err)
	}

	return u, nil
}

func (us userService) ReadList(page int, limit int) ([]User, error) {
	u := []User{}
	offset := createOffset(page, limit)
	if err := us.store.FindWithLimitAndOffset(&u, limit, offset); err != nil {
		return []User{}, err
	}

	return u, nil
}

func (us userService) Update(id string, email string) error {
	u, err := us.Read(id)
	if err != nil {
		return err
	}

	return us.store.Update(&u, map[string]interface{}{
		"email": email,
	})
}

func (us userService) UpdateWithPassword(id string, email string, hashedPassword string) error {
	u, err := us.Read(id)
	if err != nil {
		return err
	}

	return us.store.Update(&u, map[string]interface{}{
		"email":          email,
		"hashedPassword": hashedPassword,
	})
}

func (us userService) Delete(id string) error {
	u, err := us.Read(id)
	if err != nil {
		return err
	}

	return us.store.Delete(&u)
}

func createOffset(page int, limit int) int {
	if page <= 0 {
		return 0
	}

	return (page - 1) * limit
}
