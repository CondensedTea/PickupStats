package handler

import (
	"fmt"
	"net/http"

	"github.com/condensedtea/pickupRatings/v0/pkg/db"
	"github.com/labstack/echo/v4"
)

var ErrBadClass = fmt.Errorf("invalid player class: must be scout, soldier, demoman or medic")

type Handler struct {
	mongo *db.Client
}

func NewHandler(e *echo.Echo, mongo *db.Client) {
	h := &Handler{mongo: mongo}
	e.GET("/average/dpm", h.AverageDPMTotal)
	e.GET("/average/dpm/:class", h.AverageDPMOnClass)
	e.GET("/average/kdr", h.AverageKDRTotal)
	e.GET("/average/kdr/:class", h.AverageKDROnClass)
	e.GET("/average/hpm", h.AverageHealPerMin)
}

func (h *Handler) AverageDPMTotal(ctx echo.Context) error {
	results, err := h.mongo.GetAverageDPM("")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string][]db.Result{
		"stats": results,
	})
}

func (h *Handler) AverageDPMOnClass(ctx echo.Context) error {
	class := ctx.Param("class")

	if err := validateClass(class); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	results, err := h.mongo.GetAverageDPM(class)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string][]db.Result{
		"stats": results,
	})
}

func (h *Handler) AverageKDRTotal(ctx echo.Context) error {
	results, err := h.mongo.GetAverageKDR("")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string][]db.Result{
		"stats": results,
	})
}

func (h *Handler) AverageKDROnClass(ctx echo.Context) error {
	class := ctx.Param("class")

	if err := validateClass(class); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	results, err := h.mongo.GetAverageKDR(class)
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
	results, err := h.mongo.GetAverageHealsPerMin()
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
	case "scout", "soldier", "demoman", "medic":
		return nil
	default:
		return ErrBadClass
	}
}
