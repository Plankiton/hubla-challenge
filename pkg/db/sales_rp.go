package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SalesRp struct {
	PgPool *pgxpool.Pool
}

type Sale struct {
	ID      int       `json:"id"`
	Type    int       `json:"-"`
	TypeStr string    `json:"type"`
	Date    time.Time `json:"date"`
	Product string    `json:"product"`
	Value   float64   `json:"value"`
	Saler   string    `json:"saler"`
}

func SaleType(i int) string {
	return map[int]string{
		1: "Venda produtor",
		2: "Venda afiliado",
		3: "Comissão paga",
		4: "Comissão recebida",
	}[i]
}

const sqlInsertSale = `
	INSERT INTO
		sales(id, type, date, product, value, saler) VALUES (default, $1, $2, $3, $4, $5)
	RETURNING id;
	`

const sqlSelectSales = `
	SELECT id, type, date, product, value, saler FROM sales OFFSET $1 LIMIT $2
	`

const sqlCountSales = `
	SELECT count(id) FROM sales
	`

func (repo *SalesRp) Insert(ctx context.Context, sale *Sale) (*Sale, error) {
	r := repo.PgPool.QueryRow(ctx, sqlInsertSale, sale.Type, sale.Date, sale.Product, sale.Value, sale.Saler)
	err := r.Scan(&sale.ID)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (repo *SalesRp) Count(ctx context.Context) (int64, error) {
	salesCount := int64(0)

	r := repo.PgPool.QueryRow(ctx, sqlCountSales)
	err := r.Scan(&salesCount)
	if err != nil {
		return 0, err
	}

	return salesCount, nil
}

func (repo *SalesRp) Sales(ctx context.Context, o, l int) ([]*Sale, error) {
	rows, err := repo.PgPool.Query(ctx, sqlSelectSales, o, l)
	if err != nil {
		return nil, err
	}

	return repo.ScanSales(ctx, rows)
}

func (repo *SalesRp) ScanSales(ctx context.Context, row pgx.Rows) ([]*Sale, error) {
	sales := make([]*Sale, 0)

	for row.Next() {
		sale, err := repo.ScanSale(ctx, row)
		if err != nil {
			return nil, err
		}

		sales = append(sales, sale)
	}

	return sales, nil
}

func (repo *SalesRp) ScanSale(ctx context.Context, row pgx.Row) (*Sale, error) {
	sale := &Sale{}
	err := row.Scan(&sale.ID, &sale.Type, &sale.Date, &sale.Product, &sale.Value, &sale.Saler)

	sale.TypeStr = SaleType(sale.Type)
	return sale, err
}
