package domain

import "context"

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-user-manage/$GOPACKAGE -package=$GOPACKAGE -destination=store_mock.go
type Communicator interface {
	Request(ctx context.Context, to string, m map[string]string) ([]byte, error)
}
