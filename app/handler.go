package main

import (
	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(c echo.Context, err error) error {
	switch err {
	case domain.ErrNotFound, domain.ErrTokenWasExpired:
		return c.JSON(404, nil)
	case domain.ErrInvalidVerify:
		return c.JSON(401, nil)
	default:
		return c.JSON(500, nil)
	}
}
