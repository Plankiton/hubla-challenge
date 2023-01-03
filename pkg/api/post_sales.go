package api

import "github.com/labstack/echo/v4"

func (h *Handler) PostSales(c echo.Context) error {
	return c.JSON(200, echo.Map{
		"ok": true,
	})
}
