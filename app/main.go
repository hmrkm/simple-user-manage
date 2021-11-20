package main

import (
	"github.com/hmrkm/simple-user-manage/adapter"
	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/hmrkm/simple-user-manage/io"
	"github.com/hmrkm/simple-user-manage/usecase"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysqlDriver.Open(io.CreateDSN(
		config.MysqlUser,
		config.MysqlPassword,
		config.MysqlDatabase,
	)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	mysql := io.NewMysql(db)
	defer mysql.Close()

	http := io.NewHTTP(config.HttpRetryNumber, config.HttpSleepSecond)

	usd := domain.NewUserService(mysql)
	usu := usecase.NewUsers(usd, mysql)
	ua := adapter.NewUsers(usu)

	au := usecase.NewAuth(config.AuthenticationEndpoint, http, usd)
	ru := usecase.NewRights(config.RightsEndpoint, http)

	e := echo.New()
	Router(e, ua, au, ru)

	e.Logger.Fatal(e.Start(":80"))
}
