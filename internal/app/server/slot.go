package server

import (
	"database/sql"
	"hack2023/internal/app/model"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
	claims := c.Get("user").(*model.Claims)
	isKNO := claims.IsKNO

	cl, err := s.store.GetConsultationList(c.Request().Context())
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		return echo.ErrInternalServerError
	}

	tl, err := s.store.GetSlotList(c.Request().Context(), isKNO)
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		return echo.ErrInternalServerError
	}
	sl := make(map[string][]model.Slot)
	for _, p := range tl {
		if isKNO && cl[p.ID].UserID > 0 {
			cons := cl[p.ID]
			user, err := s.store.GetUserByID(c.Request().Context(), cl[p.ID].UserID)
			if err != nil {
				log.Print(err)
				return echo.ErrInternalServerError
			}
			cons.User = &user
			p.Consultation = &cons
		}
		sl[p.DateExport] = append(sl[p.DateExport], p)
	}

	return c.JSON(http.StatusOK, sl)
}
