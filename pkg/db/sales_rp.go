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

func (storage *SalesRp) Insert(ctx context.Context, sale *Sale) error {
	// var importedDataID interface{} = nil

	// _, err := storage.pgPool.Exec(ctx, sqlInsertPatient)
	// if err != nil {
	// 	return err
	// }

	return nil
}
