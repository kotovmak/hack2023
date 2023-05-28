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
	topics := make(map[int][]model.ConsultTopic)
	types := make(map[int][]model.ControlType)
	tl := &model.TypeList{}

	serviceList, err := s.store.GetServiceList(c.Request().Context())
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		return echo.ErrInternalServerError
	}
	tl.Services = append(tl.Services, serviceList...)

	consultTopics, err := s.store.GetConsultTopicList(c.Request().Context())
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, p := range consultTopics {
		topics[p.ControlTypeID] = append(topics[p.ControlTypeID], p)
	}

	controlTypes, err := s.store.GetControlTypeList(c.Request().Context())
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, p := range controlTypes {
		p.ConsultTopics = topics[p.ID]
		types[p.NadzonOrganID] = append(types[p.NadzonOrganID], p)
	}

	nadzorOrgans, err := s.store.GetNadzorOrganList(c.Request().Context())
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, p := range nadzorOrgans {
		p.ControlTypes = types[p.ID]
		tl.NadzonOrgans = append(tl.NadzonOrgans, p)
	}

	return c.JSON(http.StatusOK, tl)
}

// typelist список надзорных органов с учетом фильтрации по ключевым словам
// typelist godoc
// @Summary список надзорных органов с учетом фильтрации по ключевым словам
// @Tags consultation
// @Description список надзорных органов с учетом фильтрации по ключевым словам
// @Produce json
// @Param	text formData string true	"ключевые слова для фильтрации списка КНО" minlength(3)
// @Success 200 {object} []model.TypeList
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/typelist [post]
func (s *server) searchKNO(c echo.Context) error {
	text := c.FormValue("text")
	if len(text) < 3 {
		return c.JSON(http.StatusBadRequest, "minlenght = 3")
	}
	topics := make(map[int][]model.ConsultTopic)
	types := make(map[int][]model.ControlType)
	tl := &model.TypeList{}

	serviceList, err := s.store.GetServiceList(c.Request().Context())
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		return echo.ErrInternalServerError
	}
	tl.Services = append(tl.Services, serviceList...)

	consultTopics, err := s.store.GetConsultTopicList(c.Request().Context())
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, p := range consultTopics {
		topics[p.ControlTypeID] = append(topics[p.ControlTypeID], p)
	}

	controlTypes, err := s.store.GetControlTypeList(c.Request().Context())
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, p := range controlTypes {
		p.ConsultTopics = topics[p.ID]
		types[p.NadzonOrganID] = append(types[p.NadzonOrganID], p)
	}

	nadzorOrgans, err := s.store.GetNadzorOrganFilteredList(c.Request().Context(), text)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, p := range nadzorOrgans {
		p.ControlTypes = types[p.ID]
		tl.NadzonOrgans = append(tl.NadzonOrgans, p)
	}

	return c.JSON(http.StatusOK, tl)
}
