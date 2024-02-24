package handler

import "github.com/labstack/echo/v4"

func Debt(c echo.Context) error {
	return c.String(200, "Debt!")
}
