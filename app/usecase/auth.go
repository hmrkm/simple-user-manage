package usecase

import (
	"context"
	"encoding/json"

	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-user-manage/$GOPACKAGE -package=$GOPACKAGE -destination=auth_mock.go
type Auth interface {
	Auth(c context.Context, token string) (domain.User, error)
}

type auth struct {
	endpoint     string
	communicator domain.Communicator
	userService  domain.UserService
}

func NewAuth(ep string, c domain.Communicator, us domain.UserService) Auth {
	return auth{
		endpoint:     ep,
		communicator: c,
		userService:  us,
	}

}

type AuthResponse struct {
	User struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

func (au auth) Auth(c context.Context, token string) (domain.User, error) {
	ar := AuthResponse{}
	res, err := au.communicator.Request(c, au.endpoint, map[string]string{
		"token": token,
	})
	if err != nil {
		return domain.User{}, err
	}

	if err := json.Unmarshal(res, &ar); err != nil {
		return domain.User{}, domain.ErrInvalidValue
	}

	u, err := au.userService.Read(ar.User.Id)
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}
