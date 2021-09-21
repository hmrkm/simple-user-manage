package main

import (
	"os"

	"github.com/hmrkm/simple-user-manage/adapter"
	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/hmrkm/simple-user-manage/io"
	"github.com/hmrkm/simple-user-manage/usecase"
	"github.com/labstack/echo/v4"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	db, err := gorm.Open(mysqlDriver.Open(io.CreateDSN(mysqlUser, mysqlPassword, mysqlDatabase)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	mysql := io.Mysql{
		Conn: db,
	}
	defer mysql.Close()

	usd := domain.NewUserService(mysql)
	usu := usecase.NewUsers(usd)
	ua := adapter.NewUsers(usu)

	e := echo.New()
	Router(e, ua)

	e.Logger.Fatal(e.Start(":80"))
}
