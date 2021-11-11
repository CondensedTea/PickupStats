package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/condensedtea/PickupStats/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const defaultMinGamesAmount = 10

var ErrBadClass = fmt.Errorf("invalid player class: must be scout, soldier, demoman or medic")

type Response struct {
	Stats []db.Result `json:"stats"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

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

// AverageDPM godoc
// @Summary Player rating by average DPM.
// @Tags Ratings
// @Accept */*
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Param class path string false "Player class"
// @Param mingames path int false "Minimum games played"
// @Router /dpm [get]
func (h *Handler) AverageDPM(ctx echo.Context) error {
	class := ctx.QueryParam("class")
	minGamesRaw := ctx.QueryParam("mingames")

	minGames, err := parseMinGames(minGamesRaw)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	if err := validateClass(class); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	results, err := h.mongo.GetAverageDPM(class, minGames)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}
	return ctx.JSON(http.StatusOK, Response{Stats: results})
}

// AverageKDR godoc
// @Summary Player rating by average KDR.
// @Tags Ratings
// @Accept */*
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Param class path string false "Player class"
// @Param mingames path int false "Minimum games played"
// @Router /kdr [get]
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

// AverageHealPerMin godoc
// @Summary Medics rating by average heals given per minute.
// @Tags Ratings
// @Accept */*
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Param mingames path int false "Minimum games played"
// @Router /hpm [get]
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
