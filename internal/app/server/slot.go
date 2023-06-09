package server

import (
	"database/sql"
	"hack2023/internal/app/model"
	"log"
	"net/http"
	"time"

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
	sw := make(model.SlotWeek)
	sl := make(model.SlotList)
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
	for i, p := range sl {
		date := p[0].Date
		start := weekStartDate(date)
		start = start.AddDate(0, 0, 0)
		end := start.AddDate(0, 0, 7)
		startExport := start.Format("2 Jan")
		endExport := end.Format("2 Jan")
		s := make(model.SlotList)
		s[i] = sl[i]
		sw[startExport+" - "+endExport] = append(sw[startExport+" - "+endExport], s)
	}

	return c.JSON(http.StatusOK, sw)
}

func weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}
