package api

import (
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/plankiton/hubla-challenge/pkg/db"
)

const defaultFilename = "sales"

var saleType = map[string]string{
	"1": "Venda produtor",
	"2": "Venda afiliado",
	"3": "Comissão paga",
	"4": "Comissão recebida",
}

func (h *Handler) PostSales(c echo.Context) error {
	ctx := c.Request().Context()

	filename := defaultFilename
	if v := c.QueryParam("filename"); v != "" {
		filename = v
	}

	f, err := c.FormFile(filename)
	if err != nil {
		return c.JSON(400, echo.Map{
			"ok":  false,
			"err": fmt.Sprintf("form field of multipart needs to be \"%s\"", filename),
		})
	}

	file, err := f.Open()
	if err != nil {
		return c.JSON(400, echo.Map{
			"ok":  false,
			"err": fmt.Sprintf("file sent was invalid"),
		})
	}

	content, err := openFile(file, int(f.Size))
	if err != nil {
		return c.JSON(500, echo.Map{
			"ok":  false,
			"err": "can't read file content",
		})
	}

	var sales []*db.Sale
	for _, line := range content {
		sale := toSale(line)
		if sale != nil {
			_, err := h.rps.salesRp.Insert(ctx, sale)
			if err != nil {
				log.Printf(err.Error())
			}

			sales = append(sales, sale)
		}
	}

	return c.JSON(200, echo.Map{
		"ok":   true,
		"data": sales,
	})
}

func toSale(content []byte) *db.Sale {
	curr := 0

	typ := content[curr]
	curr++

	date := content[curr : curr+dateSize]
	curr += dateSize

	prod := string(content[curr : curr+prodSize])
	curr += prodSize

	value := content[curr : curr+valueSize]
	curr += valueSize

	saler := string(content[curr:])

	valueInt, _ := strconv.Atoi(string(value))
	valueParsed := moveFloat(valueInt, -2)

	dateParsed, _ := time.Parse(time.RFC3339, string(date))

	typParsed, _ := strconv.Atoi(string(typ))

	sale := &db.Sale{
		Type:    typParsed,
		Date:    dateParsed,
		Value:   valueParsed,
		Product: strings.TrimSpace(prod),
		Saler:   strings.TrimSpace(saler),
	}

	return sale
}

func openFile(f multipart.File, size int) ([][]byte, error) {
	content := make([]byte, size)

	_, err := f.Read(content)
	lines := [][]byte{}

	c := 0
	for i := range content {
		if content[i] == '\n' {
			if c > i {
				break
			}

			lines = append(lines, content[c:i])
			c = i + 1
		}
	}

	return lines, err
}

func moveFloat(value int, points int) float64 {
	return float64(value) * math.Pow(10.0, float64(points))
}

const (
	saleSize  = 86
	salerSize = 10
	valueSize = 10
	prodSize  = 30
	dateSize  = 25
	typSize   = 1
)
