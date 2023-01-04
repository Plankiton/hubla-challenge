package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/plankiton/hubla-challenge/pkg/db"
)

func (h *Handler) PostSales(c echo.Context) error {
	logger := c.Logger()
	ctx := c.Request().Context()

	filename := defaultFilename
	if v := c.QueryParam("filename"); v != "" {
		filename = v
	}

	f, err := c.FormFile(filename)
	if err != nil {
		logger.Errorj(log.JSON{
			"message": "wrong form field",
			"err":     err.Error(),
		})

		return c.JSON(400, echo.Map{
			"ok":  false,
			"err": fmt.Sprintf("form field of multipart needs to be \"%s\"", filename),
		})
	}

	file, err := f.Open()
	if err != nil {
		logger.Errorj(log.JSON{
			"message": "file sent was invalid",
			"err":     err.Error(),
		})

		return c.JSON(400, echo.Map{
			"ok":  false,
			"err": fmt.Sprintf("file sent was invalid"),
		})
	}

	content, err := openFile(file, int(f.Size))
	if err != nil {
		logger.Errorj(log.JSON{
			"message": "can't read file content",
			"err":     err.Error(),
		})

		return c.JSON(500, echo.Map{
			"ok":  false,
			"err": "can't read file content",
		})
	}

	var sales []*db.Sale
	for _, line := range content {
		sale := toSale(line)
		if sale != nil {
			s, err := h.rps.salesRp.Insert(ctx, sale)
			if err != nil {
				logger.Errorj(log.JSON{
					"message": "can't insert sale",
					"err":     err.Error(),
				})
			} else {
				sale = s
			}

			sales = append(sales, sale)
		}
	}

	return c.JSON(200, echo.Map{
		"ok":   true,
		"data": sales,
	})
}
