package server

import (
	"context"
	"hack2023/internal/app/model"
	"hack2023/internal/app/system"
	"log"

	"net/http"

	"firebase.google.com/go/messaging"
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

func (s *server) SendPush(ctx context.Context, mess model.PushMessage, token string) error {
	apsData := make(map[string]interface{})
	apsData["text"] = "here"
	pushMessage := &messaging.Message{
		Token: token,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Data: map[string]string{
				"apns-priority": "10",
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: mess.Title,
						Body:  mess.Body,
					},
					CustomData: apsData,
				},
			},
		},
	}
	pushMessage.Android.Notification = &messaging.AndroidNotification{
		Title: mess.Title,
		Body:  mess.Body,
	}
	_, err := s.push.Send(ctx, pushMessage)
	if err != nil {
		log.Print(err)
	}
	return nil
}
