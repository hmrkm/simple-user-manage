package main

import (
	"errors"

	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(c echo.Context, err error) error {

	if errors.As(domain.ErrNotFound, &err) ||
		errors.As(domain.ErrTokenWasExpired, &err) {
		return c.JSON(404, nil)
	} else if errors.As(domain.ErrInvalidVerify, &err) {
		return c.JSON(401, nil)
	} else if err != nil {
		return c.JSON(500, nil)
	}

	return nil
}
