package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SalesRp struct {
	PgPool *pgxpool.Pool
}

type SalesBin struct {
	Type    byte
	Date    [25]byte
	Product [30]byte
	Value   [10]byte
	Saler   [20]byte
}

type SalesModel struct {
	ID      int       `json:"id"`
	Type    string    `json:"type"`
	Date    time.Time `json:"date"`
	Product string    `json:"product"`
	Value   string    `json:"value"`
	Saler   string    `json:"saler"`
}

func (storage *SalesRp) ToModel(bin *SalesBin) *SalesModel {
	return new(SalesModel)
}

func (storage *SalesRp) Insert(ctx context.Context, sale *SalesModel) error {
	// var importedDataID interface{} = nil

	// _, err := storage.pgPool.Exec(ctx, sqlInsertPatient)
	// if err != nil {
	// 	return err
	// }

	return nil
}
