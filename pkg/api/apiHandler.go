package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/condensedtea/pickupRatings/v0/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const defaultMinGamesAmount = 10

var ErrBadClass = fmt.Errorf("invalid player class: must be scout, soldier, demoman or medic")

type Handler struct {
	mongo *db.Client
}

func NewHandler(e *echo.Echo, mongo *db.Client) {
	h := &Handler{mongo: mongo}

	api := e.Group("/api")

	api.Use(middleware.CORS())
	api.GET("/dpm", h.AverageDPM)
	api.GET("/kdr", h.AverageKDR)
	api.GET("/hpm", h.AverageHealPerMin)
	api.GET("/gamesCount", h.GamesCount)
}

func (h *Handler) AverageDPM(ctx echo.Context) error {
	class := ctx.QueryParam("class")
	minGamesRaw := ctx.QueryParam("mingames")

	minGames, err := parseMinGames(minGamesRaw)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	if err := validateClass(class); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	results, err := h.mongo.GetAverageDPM(class, minGames)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string][]db.Result{
		"stats": results,
	})
}

func (h *Handler) AverageKDR(ctx echo.Context) error {
	class := ctx.QueryParam("class")
	minGamesRaw := ctx.QueryParam("mingames")

	minGames, err := parseMinGames(minGamesRaw)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := validateClass(class); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	results, err := h.mongo.GetAverageKDR(class, minGames)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string][]db.Result{
		"stats": results,
	})
}

func (h *Handler) AverageHealPerMin(ctx echo.Context) error {
	minGamesRaw := ctx.QueryParam("mingames")

	minGames, err := parseMinGames(minGamesRaw)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	results, err := h.mongo.GetAverageHealsPerMin(minGames)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string][]db.Result{
		"stats": results,
	})
}

func (h *Handler) GamesCount(ctx echo.Context) error {
	count, err := h.mongo.GetGamesCount()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string]int64{
		"count": count / 12,
	})
}

func validateClass(class string) error {
	switch class {
	case "scout", "soldier", "demoman", "medic", "":
		return nil
	default:
		return ErrBadClass
	}
}

func parseMinGames(games string) (int, error) {
	if games == "" {
		return defaultMinGamesAmount, nil
	}
	return strconv.Atoi(games)
}
