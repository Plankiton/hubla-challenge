package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (h *Handler) GetSalesMeta(c echo.Context) error {
	logger := c.Logger()
	ctx := c.Request().Context()

	salesCount, err := h.rps.salesRp.Count(ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"ok":  false,
			"err": "error counting sales",
		})
	}

	salesSum, err := h.rps.salesRp.Sum(ctx)
	if err != nil {
		logger.Errorj(log.JSON{
			"message": "error summing sale",
			"err":     err.Error(),
		})

		return c.JSON(500, echo.Map{
			"ok":  false,
			"err": "error summing sales",
		})
	}

	return c.JSON(200, echo.Map{
		"ok": true,
		"meta": echo.Map{
			"total": salesSum,
			"count": salesCount,
		},
	})
}
