package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/condensedtea/pickupRatings/v0/pkg/db"
	"github.com/labstack/echo/v4"
)

const defaultMinGamesAmount = 10

var ErrBadClass = fmt.Errorf("invalid player class: must be scout, soldier, demoman or medic")

type Handler struct {
	mongo *db.Client
}

func NewHandler(e *echo.Echo, mongo *db.Client) {
	h := &Handler{mongo: mongo}

	e.GET("/average/dpm", h.AverageDPM)
	e.GET("/average/kdr", h.AverageKDR)
	e.GET("/average/hpm", h.AverageHealPerMin)
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
	class := ctx.Param("class")
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
