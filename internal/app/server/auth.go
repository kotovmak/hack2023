package server

import (
	"hack2023/internal/app/model"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// login Получение токена авторизации
// login godoc
// @Summary Получение токена авторизации
// @Tags auth
// @Description Получение токена авторизации
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
		return echo.ErrBadRequest
	}

	user, err := s.store.GetUserByLogin(c.Request().Context(), login)
	if err != nil {
		log.Print(err)
		return echo.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Print(err)
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &model.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 1)},
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
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 1)},
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
