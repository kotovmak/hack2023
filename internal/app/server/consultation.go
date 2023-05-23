package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// getTypeList список надзорных органов, видов контроля, тем консультирования
// getTypeList godoc
// @Summary список надзорных органов, видов контроля, тем консультирования
// @Tags consultation
// @Description список надзорных органов, видов контроля, тем консультирования
// @Produce json
// @Success 200 {object} []model.TypeList
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/typelist [get]
func (s *server) getTypeList(c echo.Context) error {

	tl, err := s.store.GetTypeList(c.Request().Context())
	if err != nil {
		log.Print(err)
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tl)
}

// getSlotList список доступных слотов
// getSlotList godoc
// @Summary список доступных слотов
// @Tags consultation
// @Description список доступных слотов
// @Produce json
// @Success 200 {object} []model.Slot
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/slot [get]
func (s *server) getSlotList(c echo.Context) error {

	tl, err := s.store.GetSlotList(c.Request().Context())
	if err != nil {
		log.Print(err)
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tl)
}
