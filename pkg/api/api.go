package api

import (
	"math"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/plankiton/hubla-challenge/pkg/db"
)

type Repositories struct {
	salesRp *db.SalesRp
	Close   func()
}

type Handler struct {
	rps *Repositories
}

func (h Handler) Close() {
	h.rps.Close()
}

func New(rps *Repositories) Handler {
	return Handler{
		rps,
	}
}

func NewRepositories(pgPool *pgxpool.Pool) *Repositories {
	return &Repositories{
		salesRp: &db.SalesRp{
			pgPool,
		},
		Close: func() {
			pgPool.Close()
		},
	}
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

const defaultFilename = "sales"
const defaultLimit = 100
const defaultOffset = 0

var saleType = map[string]string{
	"1": "Venda produtor",
	"2": "Venda afiliado",
	"3": "Comissão paga",
	"4": "Comissão recebida",
}
