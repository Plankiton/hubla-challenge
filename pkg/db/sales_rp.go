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
	ID           int       `json:"id"`
	Type         int       `json:"-"`
	TypeStr      string    `json:"type"`
	Date         time.Time `json:"date"`
	Product      string    `json:"product"`
	Value        float64   `json:"-"`
	DisplayValue float64   `json:"value"`
	Seller       string    `json:"seller"`
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
		sales(id, type, date, product, value, seller) VALUES (default, $1, $2, $3, $4, $5)
	RETURNING id;
	`

const sqlSelectSales = `
	SELECT id, type, date, product, value, seller FROM sales OFFSET $1 LIMIT $2
	`

const sqlCountSales = `
	SELECT count(id) FROM sales
	`

const sqlSumSales = `
	SELECT sum(value) FROM sales
	`

func (repo *SalesRp) Insert(ctx context.Context, sale *Sale) (*Sale, error) {
	r := repo.PgPool.QueryRow(ctx, sqlInsertSale, sale.Type, sale.Date, sale.Product, sale.Value, sale.Seller)
	err := r.Scan(&sale.ID)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (repo *SalesRp) Sum(ctx context.Context) (float64, error) {
	salesSum := 0.0

	count, err := repo.Count(ctx)
	if err != nil {
		return 0.0, err
	}

	if count > 0 {
		r := repo.PgPool.QueryRow(ctx, sqlSumSales)
		err := r.Scan(&salesSum)
		if err != nil && err != pgx.ErrNoRows {
			return 0.0, err
		}
	}

	return salesSum, nil
}

func (repo *SalesRp) Count(ctx context.Context) (int64, error) {
	salesCount := int64(0)

	r := repo.PgPool.QueryRow(ctx, sqlCountSales)
	err := r.Scan(&salesCount)
	if err != nil && err != pgx.ErrNoRows {
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
	err := row.Scan(&sale.ID, &sale.Type, &sale.Date, &sale.Product, &sale.Value, &sale.Seller)

	sale.TypeStr = SaleType(sale.Type)
	sale.DisplayValue = sale.Value
	if sale.Value < 0 {
		sale.DisplayValue *= -1
	}

	return sale, err
}
