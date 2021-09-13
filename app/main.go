package main

import (
	"os"

	"github.com/hmrkm/simple-user-manage/adapter"
	"github.com/hmrkm/simple-user-manage/io"
	"github.com/hmrkm/simple-user-manage/usecase"
	"github.com/labstack/echo/v4"
)

func main() {

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	mysql := io.OpenMysql(mysqlUser, mysqlPassword, mysqlDatabase)
	defer mysql.Close()

	usu := usecase.NewUserService(mysql)
	ua := adapter.NewUsersAdapter(usu)

	e := echo.New()
	e.GET("/v1/users", func(c echo.Context) error {
		req := adapter.GetV1UsersParams{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, nil)
		}

		res, err := ua.GetList(req)

		if err != nil {
			return c.JSON(403, nil)
		}

		return c.JSON(200, res)
	})

	e.Logger.Fatal(e.Start(":80"))
}
