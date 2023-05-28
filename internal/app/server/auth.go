package server

import (
	"database/sql"
	"hack2023/internal/app/model"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// getUser Получение данных о текущем авторизованном пользователе
// getUser godoc
// @Summary Получение данных о текущем авторизованном пользователе
// @Tags auth
// @Description Для представителя КНО есть дополнительные поля, необходимые для сложебных методов: Подтверждение консультации
// @Produce json
// @Success 200 {object} []model.Account
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/user [get]
func (s *server) getUser(c echo.Context) error {
	claims := c.Get("user").(*model.Claims)
	login := claims.Login
	tl, err := s.store.GetUserByLogin(c.Request().Context(), login)
	if err != nil {
		log.Print(err)
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tl)
}

// saveAppToken Сохранение токена приложения пользователя
// saveAppToken godoc
// @Summary Сохранение токена приложения пользователя
// @Tags auth
// @Description Необходимо для дальнейшей рассылки push сообщений на телефоны
// @Produce json
// @Success 201 {object} []model.Token
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/user [post]
func (s *server) saveAppToken(c echo.Context) error {
	claims := c.Get("user").(*model.Claims)
	userID := claims.ID
	token := c.FormValue("token")
	if len(token) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "token is required")
	}

	tl := model.Token{UserID: userID, Token: token}
	err := s.store.DeleteAppToken(c.Request().Context(), token)
	if err != nil {
		log.Print(err)
		return echo.ErrInternalServerError
	}
	tl, err = s.store.SaveAppToken(c.Request().Context(), tl)
	if err != nil {
		log.Print(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, tl)
}

// login Получение токена авторизации
// login godoc
// @Summary Получение токена авторизации
// @Tags auth
// @Description Необходим для использования во всех остальных запросах
// @Produce json
// @Param login formData string true "login"
// @Param password formData string true "password"
// @Success 200 {object} []model.AuthResponse
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/login [post]
func (s *server) login(c echo.Context) error {
	login := c.FormValue("login")
	password := c.FormValue("password")

	// Throws unauthorized error
	if len(login) == 0 || len(password) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "login and password is required fields")
	}

	user, err := s.store.GetUserByLogin(c.Request().Context(), login)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Set custom claims
	claims := &model.Claims{
		ID:            user.ID,
		Login:         user.Login,
		Email:         user.Email,
		Name:          user.Name,
		IsKNO:         user.IsKNO,
		NadzonOrganID: user.NadzonOrganID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * 30)},
		},
	}

	refreshClaims := &model.RefreshClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * 30)},
		},
	}

	// Create token with claims
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Generate encoded token and send it as response.
	at, err := accessToken.SignedString([]byte(s.config.SigningKey))
	if err != nil {
		log.Print(err)
		return err
	}

	rt, err := refreshToken.SignedString([]byte(s.config.SigningKey))
	if err != nil {
		log.Print(err)
		return err
	}

	err = s.store.DeleteRefreshToken(c.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	err = s.store.SaveRefreshToken(c.Request().Context(), user.ID, rt)
	if err != nil {
		log.Print(err)
		return err
	}

	return c.JSON(http.StatusOK, model.AuthResponse{
		AccessToken:  at,
		RefreshToken: rt,
	})
}

// handleToken Получение токена авторизации по refresh токену
// handleToken godoc
// @Summary Получение токена авторизации по refresh токену
// @Tags auth
// @Description Получение токена авторизации по refresh токену
// @Produce json
// @Param refresh_token formData string true "refresh_token"
// @Success 200 {object} []model.AuthResponse
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/token [post]
func (s *server) handleToken(c echo.Context) error {
	rtOld := c.FormValue("refresh_token")

	if len(rtOld) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errInvalidToken.Error())
	}

	token, err := jwt.ParseWithClaims(rtOld, &model.RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errWrongSingingMethod
		}
		return []byte(s.config.SigningKey), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	parsedClaims, ok := token.Claims.(*model.RefreshClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, errInvalidToken.Error())
	}

	user, err := s.store.GetUserByRefreshToken(c.Request().Context(), rtOld)
	if err != nil {
		log.Print(err)
		return echo.ErrUnauthorized
	}

	if parsedClaims.RegisteredClaims.ExpiresAt.Time.Unix() < time.Now().Unix() {
		err = s.store.DeleteRefreshToken(c.Request().Context(), user.ID)
		log.Print(err)
		if err != nil {
			return err
		}
		return echo.ErrUnauthorized
	}

	// Set custom claims
	newClaims := &model.Claims{
		ID:            user.ID,
		Login:         user.Login,
		Email:         user.Email,
		Name:          user.Name,
		IsKNO:         user.IsKNO,
		NadzonOrganID: user.NadzonOrganID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * 30)},
		},
	}

	refreshClaims := &model.RefreshClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * 30)},
		},
	}

	// Create token with claims
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Generate encoded token and send it as response.
	at, err := accessToken.SignedString([]byte(s.config.SigningKey))
	if err != nil {
		log.Print(err)
		return err
	}

	rt, err := refreshToken.SignedString([]byte(s.config.SigningKey))
	if err != nil {
		log.Print(err)
		return err
	}

	err = s.store.DeleteRefreshToken(c.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	err = s.store.SaveRefreshToken(c.Request().Context(), user.ID, rt)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.AuthResponse{
		AccessToken:  at,
		RefreshToken: rt,
	})
}
