package main

import (
	"context"
	"log"

	"github.com/condensedtea/PickupStats/v0/docs"
	"github.com/condensedtea/PickupStats/v0/pkg/api"
	"github.com/condensedtea/PickupStats/v0/pkg/config"
	"github.com/condensedtea/PickupStats/v0/pkg/db"
	"github.com/condensedtea/PickupStats/v0/pkg/frontend"
	"github.com/condensedtea/PickupStats/v0/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/condensedtea/PickupStats/v0/docs"
)

const configPath = "config.yaml"

const loglevel = "debug"

var Version = "dev"

// @title Pickup Stats API
// @description API for pickup stats collected with LogWatcher.

// @BasePath /api
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

	client, err := db.NewClient(ctx, cfg.DSN, cfg.Database, cfg.GameCollection, cfg.NameCollection)
	if err != nil {
		l.Fatalf("Failed to conntect to mongodb: %v", err)
	}

	api.NewHandler(e, client)
	frontend.NewHandler(e)

	docs.SwaggerInfo.Version = Version
	e.Use(middleware.Recover())
	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
