package server

import (
	"database/sql"
	"hack2023/internal/app/model"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	BOTID = 0
)

// getMessageList список сообщения чат-бота текущего пользователя
// getMessageList godoc
// @Summary список сообщения чат-бота текущего пользователя
// @Tags consultation
// @Description список сообщения чат-бота текущего пользователя
// @Produce json
// @Success 200 {object} model.Message
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/chat [get]
func (s *server) getMessageList(c echo.Context) error {
	claims := c.Get("user").(*model.Claims)

	tl, err := s.store.GetMessagesList(c.Request().Context(), claims.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		return echo.ErrInternalServerError
	}
	if len(tl) == 0 {
		date := time.Now()
		dateExport := date.Format("2006-01-02")
		m, err := s.store.AddMessage(c.Request().Context(), model.Message{
			UserID:     claims.ID,
			SendByID:   BOTID,
			Date:       date,
			DateExport: dateExport,
			Text: `Приветствую, ` + claims.Name +
				`На связи чат-бот.`,
		})
		if err != nil {
			return echo.ErrInternalServerError
		}
		tl = append(tl, m)
		m, err = s.store.AddMessage(c.Request().Context(), model.Message{
			UserID:     claims.ID,
			SendByID:   BOTID,
			Date:       date,
			DateExport: dateExport,
			Text: `Я здесь, чтобы сэкономить ваше
время.`,
		})
		if err != nil {
			return echo.ErrInternalServerError
		}
		tl = append(tl, m)
		m, err = s.store.AddMessage(c.Request().Context(), model.Message{
			UserID:     claims.ID,
			SendByID:   BOTID,
			Date:       date,
			DateExport: dateExport,
			Text: `Собрал ответы на популярные 
вопросы органам контроля и информационные материалы.`,
		})
		if err != nil {
			return echo.ErrInternalServerError
		}
		tl = append(tl, m)
	}

	cl := make(map[string][]model.Message)
	for _, t := range tl {
		cl[t.DateExport] = append(cl[t.DateExport], t)
	}

	return c.JSON(http.StatusOK, cl)
}

// addMessage отправить сообщение чат-боту и получить ответ
// addMessage godoc
// @Summary отправить сообщение чат-боту и получить ответ
// @Tags consultation
// @Description отправить сообщение чат-боту и получить ответ
// @Produce json
// @Param	text	formData string true	"текст сообщения" minlength(1)
// @Success 202 {object} model.Message
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/chat [post]
func (s *server) addMessage(c echo.Context) error {
	claims := c.Get("user").(*model.Claims)

	cl := model.Message{}
	if err := c.Bind(&cl); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cl.UserID = claims.ID
	cl.SendByID = claims.ID
	cl.Date = time.Now()
	cl.DateExport = cl.Date.Format("2006-01-02")
	if err := c.Validate(&cl); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cl, err := s.store.AddMessage(c.Request().Context(), cl)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//TODO: поиск по ключевым словам

	returnMessage := model.Message{
		UserID:     cl.UserID,
		SendByID:   BOTID,
		Date:       time.Now().Add(time.Minute * 1),
		DateExport: cl.DateExport,
		Text:       "тут будет ответ бота",
	}
	cl, err = s.store.AddMessage(c.Request().Context(), returnMessage)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, cl)
}

// getButtonList список кнопок чата
// getButtonList godoc
// @Summary список кнопок чата
// @Tags consultation
// @Description список кнопок чата
// @Produce json
// @Success 200 {object} model.Button
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Security ApiKeyAuth
// @Router /v1/button [get]
func (s *server) getButtonList(c echo.Context) error {
	tl, err := s.store.GetButtonList(c.Request().Context())
	if err != nil {
		log.Print(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tl)
}
