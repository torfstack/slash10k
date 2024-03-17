package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	frontend "scurvy10k/templ"
)

func ServeFrontend(c echo.Context) error {
	return frontend.
		Debt().
		Render(context.Background(), c.Response())
}
