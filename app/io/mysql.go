package io

import (
	"errors"
	"fmt"

	"github.com/hmrkm/simple-user-manage/usecase"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Mysql struct {
	Conn *gorm.DB
}

func OpenMysql(user string, password string, database string) Mysql {
	dsn := fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	return Mysql{
		Conn: db,
	}
}

func (m Mysql) Close() {
	db, err := m.Conn.DB()
	if err != nil {
		panic(err)
	}

	db.Close()
}

func (m Mysql) Find(destAddr interface{}, conds string, params ...interface{}) error {
	return m.Conn.Find(destAddr, conds, params).Error
}

func (m Mysql) First(destAddr interface{}, conds string, params ...interface{}) error {
	err := m.Conn.First(destAddr, conds, params).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usecase.ErrNotFound
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
