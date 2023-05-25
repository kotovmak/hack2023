package server

import (
	"database/sql"
	"hack2023/internal/app/model"
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

// getConsultationList список активных и завершенных консультаций
// getConsultationList godoc
// @Summary список активных и завершенных консультаций
// @Tags consultation
// @Description список активных и завершенных консультаций
// @Produce json
// @Success 200 {object} model.Consultations
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/consultation [get]
func (s *server) getConsultationList(c echo.Context) error {

	tl, err := s.store.GetConsultationList(c.Request().Context())
	if err != nil {
		log.Print(err)
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tl)
}

// addConsultation записаться на консультацию
// addConsultation godoc
// @Summary записаться на консультацию
// @Tags consultation
// @Description записаться на консультацию
// @Produce json
// @Param	nadzor_organ_id	formData int true	"id надзорного органа" minimum(1)
// @Param	control_type_id	formData int true	"id типа контроля" minimum(1)
// @Param	consult_topic_id formData int true "id темы консультации" minimum(1)
// @Param	user_id	formData int true	"id пользователя" minimum(1)
// @Param	time formData string true	"время в формате '03:00'"
// @Param	date formData string true	"дата в формате '2006-02-01'"
// @Param	question formData string true	"вопрос в свободной форме"
// @Param	is_need_letter formData bool true	"нужно ли письменное разъяснение"
// @Success 201 {object} model.Consultation
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/consultation [post]
func (s *server) addConsultation(c echo.Context) error {
	cl := model.Consultation{}
	if err := c.Bind(&cl); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err := s.store.AddConsultation(c.Request().Context(), cl)
	if err != nil {
		log.Print(err)
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	cl.DateExport = cl.Date.Format("2006-01-02")

	return c.JSON(http.StatusCreated, cl)
}
