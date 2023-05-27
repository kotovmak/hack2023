package server

import (
	"database/sql"
	"hack2023/internal/app/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

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
		return echo.ErrInternalServerError
	}

	cl := &model.Consultations{}
	for _, p := range tl {
		if p.Date.Unix() > time.Now().Unix() {
			cl.Active = append(cl.Active, p)
		} else {
			cl.Finished = append(cl.Finished, p)
		}
	}

	return c.JSON(http.StatusOK, cl)
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
// @Param	slot_id	formData int true	"id слота с временем и датой консультации" minimum(1)
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
	claims := c.Get("user").(*model.Claims)

	cl := model.Consultation{}
	if err := c.Bind(&cl); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if claims.ID > 0 {
		cl.UserID = claims.ID
	}
	cl.VKSLink = "https://peregovorka.mos.ru/" + uuid.New().String()
	if err := c.Validate(&cl); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	slot, err := s.store.GetSlot(c.Request().Context(), cl.SlotID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, errSlotBusy.Error())
		}
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	cl, err = s.store.AddConsultation(c.Request().Context(), cl)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = s.store.CloseSlot(c.Request().Context(), slot.ID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cl.DateExport = cl.Date.Format("2006-01-02")
	return c.JSON(http.StatusCreated, cl)
}

// deleteConsultation Отменить запись на консультацию
// deleteConsultation godoc
// @Summary Отменить запись на консультацию
// @Tags consultation
// @Description Отменить запись на консультацию
// @Produce json
// @Param	id	query int true	"id консультации которую нужно отменить" minimum(1)
// @Success 200 {object} model.Consultation
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/consultation [delete]
func (s *server) deleteConsultation(c echo.Context) error {
	consultationID := c.QueryParam("id")
	if len(consultationID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}
	cl, err := s.store.GetConsultation(c.Request().Context(), consultationID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = s.store.DeleteConsultation(c.Request().Context(), consultationID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	cl.IsDeleted = true
	err = s.store.OpenSlot(c.Request().Context(), cl.SlotID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cl.DateExport = cl.Date.Format("2006-01-02")
	return c.JSON(http.StatusOK, cl)
}

// applyConsultation Подтверждение консультации со стороны КНО
// applyConsultation godoc
// @Summary Подтверждение консультации со стороны КНО
// @Tags consultation
// @Description Подтверждение консультации со стороны КНО
// @Produce json
// @Param	id formData int true	"id консультации которую нужно подтвердить" minimum(1)
// @Param	apply formData bool true	"Подтвердить или нет консультацию"
// @Success 201 {object} model.Consultation
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/consultation [patch]
func (s *server) applyConsultation(c echo.Context) error {
	claims := c.Get("user").(*model.Claims)
	isKNO := claims.IsKNO

	if !isKNO {
		return echo.NewHTTPError(http.StatusNotAcceptable, errOnlyKNO.Error())
	}

	consultationID := c.FormValue("id")
	if len(consultationID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	apply, err := strconv.ParseBool(c.FormValue("apply"))
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cl, err := s.store.GetConsultation(c.Request().Context(), consultationID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if apply {
		err = s.store.ApplyConsultation(c.Request().Context(), consultationID)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		cl.IsConfirmed = true
		return c.JSON(http.StatusOK, cl)
	}

	err = s.store.DeleteConsultation(c.Request().Context(), consultationID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	cl.IsDeleted = true
	err = s.store.OpenSlot(c.Request().Context(), cl.SlotID)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, cl)
}
