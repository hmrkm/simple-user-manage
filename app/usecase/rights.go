package usecase

import (
	"context"

	"github.com/hmrkm/simple-user-manage/domain"
)

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-user-manage/$GOPACKAGE -package=$GOPACKAGE -destination=rights_mock.go
type Rights interface {
	Verify(c context.Context, userID string, resource string) error
}

type rights struct {
	endpoint     string
	communicator domain.Communicator
}

func NewRights(ed string, c domain.Communicator) Rights {
	return rights{
		endpoint:     ed,
		communicator: c,
	}
}

func (ru rights) Verify(c context.Context, userID string, resource string) error {
	_, err := ru.communicator.Request(c, ru.endpoint, map[string]string{
		"user_id":  userID,
		"resource": resource,
	})
	if err != nil {
		return err
	}

	return nil
}
