package domain

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-user-manage/$GOPACKAGE -package=$GOPACKAGE -destination=store_mock.go
type Store interface {
	Find(destAddr interface{}, cond string, params ...interface{}) error
	First(destAddr interface{}, cond string, params ...interface{}) error
	FindWithLimitAndOffset(destAddr interface{}, limit int, offset int) error
	Create(targetAddr interface{}) error
	Update(targetAddr interface{}, params map[string]interface{}) error
	Delete(targetAddr interface{}) error
	Count(targetAddr interface{}, count *int64) error
	IsNotFoundError(error) bool
}
