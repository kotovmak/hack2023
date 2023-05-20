package server

import (
	"hack2023/internal/app/model"
	"hack2023/internal/app/system"

	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *server) ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func (s *server) handleStatus(ctx echo.Context) error {
	err := system.Healthz()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, model.Status{
		Status: model.StatusOK,
	})
}

func (s *server) handleVersion(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.Version{
		Tier:     s.config.Tier,
		Version:  s.config.Version,
		Revision: s.config.Revision,
	})
}
