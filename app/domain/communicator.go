package domain

import "context"

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-user-manage/$GOPACKAGE -package=$GOPACKAGE -destination=communicator_mock.go
type Communicator interface {
	Request(ctx context.Context, to string, body interface{}) ([]byte, error)
}
