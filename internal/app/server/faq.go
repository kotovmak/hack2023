package server

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// getFAQList список вопросов и ответов
// getFAQList godoc
// @Summary список вопросов и ответов
// @Tags consultation
// @Description список вопросов и ответов
// @Produce json
// @Success 200 {object} []model.FAQ
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/faq [get]
func (s *server) getFAQList(c echo.Context) error {

	tl, err := s.store.GetFAQList(c.Request().Context())
	if err != nil {
		log.Print(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tl)
}

// searchFAQ список вопросов и ответов c фтльтрацией
// searchFAQ godoc
// @Summary список вопросов и ответов c фтльтрацией
// @Tags consultation
// @Description список вопросов и ответов c фтльтрацией
// @Produce json
// @Param	text formData string true	"ключевые слова для фильтрации списка КНО" minlength(3)
// @Success 200 {object} []model.FAQ
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/faq [post]
func (s *server) searchFAQ(c echo.Context) error {
	text := c.FormValue("text")
	if len(text) < 3 {
		return c.JSON(http.StatusBadRequest, "minlenght = 3")
	}
	tl, err := s.store.SearchFAQList(c.Request().Context(), text)
	if err != nil {
		log.Print(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tl)
}
