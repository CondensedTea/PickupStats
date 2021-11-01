package main

import (
	"context"
	"log"

	"github.com/condensedtea/pickupRatings/v0/pkg/api"
	"github.com/condensedtea/pickupRatings/v0/pkg/config"
	"github.com/condensedtea/pickupRatings/v0/pkg/db"
	"github.com/condensedtea/pickupRatings/v0/pkg/frontend"
	"github.com/condensedtea/pickupRatings/v0/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const configPath = "config.yaml"

const loglevel = "debug"

func main() {
	e := echo.New()
	ctx := context.Background()

	l, err := logger.SetLogger(e, loglevel)
	if err != nil {
		log.Fatalf("Failed to parse loglevel: %v", err)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		l.Fatalf("Failed to parse config: %v", err)
	}

	client, err := db.NewClient(ctx, cfg.DSN, cfg.Database, cfg.Collection)
	if err != nil {
		l.Fatalf("Failed to conntect to mongodb: %v", err)
	}

	api.NewHandler(e, client)
	frontend.NewHandler(e)

	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":1323"))
}
