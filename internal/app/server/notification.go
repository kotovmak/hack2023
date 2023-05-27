package server

import (
	"hack2023/internal/app/model"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// getNotificationList список уведомлений текущего пользователя
// getNotificationList godoc
// @Summary список уведомлений текущего пользователя
// @Tags consultation
// @Description список уведомлений текущего пользователя
// @Produce json
// @Success 200 {object} model.Notification
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/notification [get]
func (s *server) getNotificationList(c echo.Context) error {
	claims := c.Get("user").(*model.Claims)

	tl, err := s.store.GetNotificationList(c.Request().Context(), claims.ID)
	if err != nil {
		log.Print(err)
		return echo.ErrInternalServerError
	}
	cl := make(map[string][]model.Notification)
	for _, t := range tl {
		cl[t.DateExport] = append(cl[t.DateExport], t)
	}

	return c.JSON(http.StatusOK, cl)
}
