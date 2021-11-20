package main

import (
	"github.com/hmrkm/simple-user-manage/usecase"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(au usecase.Auth) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("auth_token")
			user, err := au.Auth(c.Request().Context(), token)
			if err != nil {
				return c.JSON(401, nil)
			}

			c.Set("authorization_user_id", user.Id)

			return next(c)
		}
	}
}

func RightsMiddleware(ru usecase.Rights) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("authorization_user_id").(string)
			err := ru.Verify(c.Request().Context(), userID, c.Request().URL.Path)
			if err != nil {
				return c.JSON(403, nil)
			}

			return next(c)
		}
	}
}
