package usecase

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-user-manage/$GOPACKAGE -package=$GOPACKAGE -destination=store_mock.go
type Store interface {
	Find(destAddr interface{}, conds string, params ...interface{}) error
	First(destAddr interface{}, conds string, params ...interface{}) error
	FindWithLimitAndOffset(destAddr interface{}, limit int, offset int) error
	Create(targetAddr interface{}) error
	Update(targetAddr interface{}, params map[string]interface{}) error
	Delete(targetAddr interface{}) error
}
