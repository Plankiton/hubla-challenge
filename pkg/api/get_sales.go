package api

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetSales(c echo.Context) error {
	ctx := c.Request().Context()

	limit := defaultLimit
	if l, _ := strconv.Atoi(c.QueryParam("limit")); l > 0 {
		limit = l
	}

	offset := defaultOffset
	if o, _ := strconv.Atoi(c.QueryParam("offset")); o > 0 {
		offset = o
	}

	sales, err := h.rps.salesRp.Sales(ctx, offset, limit)
	if err != nil {
		return c.JSON(500, echo.Map{
			"ok":  false,
			"err": "error searching for sales",
		})
	}

	salesCount, err := h.rps.salesRp.Count(ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"ok":  false,
			"err": "error counting sales",
		})
	}

	return c.JSON(200, echo.Map{
		"ok":   true,
		"data": sales,
		"meta": echo.Map{
			"offset": offset,
			"limit":  limit,
			"total":  salesCount,
		},
	})
}
