package io

import (
	"errors"
	"fmt"

	"github.com/hmrkm/simple-user-manage/domain"

	"gorm.io/gorm"
)

type Mysql struct {
	Conn GormConn
}

func CreateDSN(user string, password string, database string) string {
	return fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, database)
}

func (m Mysql) Close() error {
	db, err := m.Conn.DB()
	if err != nil {
		return err
	}

	db.Close()

	return nil
}

func (m Mysql) Find(destAddr interface{}, conds string, params ...interface{}) error {
	return m.Conn.Find(destAddr, conds, params).Error
}

func (m Mysql) First(destAddr interface{}, conds string, params ...interface{}) error {
	err := m.Conn.First(destAddr, conds, params).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}

	return err
}

func (m Mysql) FindWithLimitAndOffset(destAddr interface{}, limit int, offset int) error {
	return m.Conn.Limit(limit).Offset(offset).Find(destAddr).Error
}

func (m Mysql) Create(targetAddr interface{}) error {
	return m.Conn.Create(targetAddr).Error
}

func (m Mysql) Update(targetAddr interface{}, params map[string]interface{}) error {
	return m.Conn.Model(targetAddr).Updates(params).Error
}

func (m Mysql) Delete(targetAddr interface{}) error {
	return m.Conn.Delete(targetAddr).Error
}

func (m Mysql) Count(targetAddr interface{}, count *int64) error {
	return m.Conn.Model(targetAddr).Count(count).Error
}

func (m Mysql) IsNotFoundError(err error) bool {
	return errors.Is(gorm.ErrRecordNotFound, err)
}
