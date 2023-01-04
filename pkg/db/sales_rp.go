package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SalesRp struct {
	PgPool *pgxpool.Pool
}

type Sale struct {
	ID      int       `json:"id"`
	Type    int       `json:"type"`
	Date    time.Time `json:"date"`
	Product string    `json:"product"`
	Value   float64   `json:"value"`
	Saler   string    `json:"saler"`
}

const sqlInsertSale = `
	INSERT INTO
		sales(id, type, date, product, value, saler) VALUES (default, $1, $2, $3, $4, $5)
	RETURNING id;
	`

func (storage *SalesRp) Insert(ctx context.Context, sale *Sale) (*Sale, error) {
	r := storage.PgPool.QueryRow(ctx, sqlInsertSale, sale.Type, sale.Date, sale.Product, sale.Value, sale.Saler)
	err := r.Scan(&sale.ID)
	if err != nil {
		return nil, err
	}

	return sale, nil
}
