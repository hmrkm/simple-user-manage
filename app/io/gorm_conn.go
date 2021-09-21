package io

import (
	"database/sql"

	"gorm.io/gorm"
)

//go:generate mockgen -source=$GOFILE -self_package=github.com/hmrkm/simple-auth/$GOPACKAGE -package=$GOPACKAGE -destination=gorm_conn_mock.go
type GormConn interface {
	DB() (*sql.DB, error)
	Find(destAddr interface{}, conds ...interface{}) *gorm.DB
	First(destAddr interface{}, conds ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Updates(values interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Limit(limit int) *gorm.DB
	Offset(offset int) *gorm.DB
	Model(targetAddr interface{}) *gorm.DB
}
