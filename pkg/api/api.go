package api

import (
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
