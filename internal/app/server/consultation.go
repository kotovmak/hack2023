package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// getTypeList Получить список справочников
// getTypeList godoc
// @Summary Получить список справочников
// @Tags contact
// @Description Получить список справочников
// @Produce json
// @Success 200 {object} []model.TypeList
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/getTypeList [get]
func (s *server) getTypeList(c echo.Context) error {
	inn := c.FormValue("inn")

	if len(inn) < 10 || len(inn) > 12 || len(inn) == 11 {
		return echo.ErrBadRequest
	}

	contacts, err := s.store.GetTypeList(c.Request().Context(), inn)
	if err != nil {
		log.Print(err)
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, contacts)
}
