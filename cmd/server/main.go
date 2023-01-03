package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/plankiton/hubla-challenge/cmd/server/internal/config"
	"github.com/plankiton/hubla-challenge/pkg/api"
)

func main() {
	ctx := context.Background()
	server, err := createEchoServer(ctx, config.New())
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}

func createRepositories(ctx context.Context, config config.Config) (*api.Repositories, error) {
	conf, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse postgres config: %w", err)
	}

	pgPool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	return api.NewRepositories(pgPool), nil
}

func createEchoServer(ctx context.Context, config config.Config) (*http.Server, error) {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Static("/", config.StaticPath)

	repos, err := createRepositories(ctx, config)
	if err != nil {
		return nil, err
	}

	handler := api.New(repos)

	apiG := e.Group("/api")
	apiG.POST("/sales", handler.PostSales)

	return &http.Server{
		Handler: e,
		Addr:    ":2345",
	}, nil
}
