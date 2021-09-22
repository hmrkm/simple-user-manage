package main

import (
	"github.com/hmrkm/simple-user-manage/adapter"
	"github.com/hmrkm/simple-user-manage/usecase"
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo, ua adapter.Users, au usecase.Auth) {
	g := e.Group("/v1/users")
	g.Use(AuthMiddleware(au))
	g.POST("/list", func(c echo.Context) error {
		req := adapter.RequestUsersList{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, nil)
		}

		res, err := ua.List(req)

		if err != nil {
			return ErrorHandler(c, err)
		} else {
			return c.JSON(200, res)
		}
	})

	g.POST("/create", func(c echo.Context) error {
		req := adapter.RequestUsersCreate{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, nil)
		}

		err := ua.Create(req)

		if err != nil {
			return ErrorHandler(c, err)
		} else {
			return c.JSON(200, nil)
		}
	})

	g.POST("/detail", func(c echo.Context) error {
		req := adapter.RequestUsersDetail{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, nil)
		}

		res, err := ua.Detail(req)

		if err != nil {
			return ErrorHandler(c, err)
		} else {
			return c.JSON(200, res)
		}

	})

	g.POST("/update", func(c echo.Context) error {
		req := adapter.RequestUsersUpdate{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, nil)
		}

		err := ua.Update(req)

		if err != nil {
			return ErrorHandler(c, err)
		} else {
			return c.JSON(200, nil)
		}
	})

	g.POST("/delete", func(c echo.Context) error {
		req := adapter.RequestUsersDelete{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, nil)
		}

		err := ua.Delete(req)

		if err != nil {
			return ErrorHandler(c, err)
		} else {
			return c.JSON(200, nil)
		}
	})
}
