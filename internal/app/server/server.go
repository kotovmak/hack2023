package server

import (
	"fmt"
	"hack2023/internal/app/config"
	"hack2023/internal/app/model"
	"hack2023/internal/app/store"

	_ "hack2023/docs"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type server struct {
	router *echo.Echo
	store  *store.Store
	config config.Config
}

// NewServer инициализируем сервер
func NewServer(store *store.Store, config config.Config) *server {
	s := &server{
		router: echo.New(),
		store:  store,
		config: config,
	}

	// Конфигурируем роутинг
	s.configureRouter()
	return s
}

// Start Включаем прослушивание HTTP протокола
func (s *server) Start(config config.Config) error {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err := s.router.Start(address)
	if err != nil {
		return err
	}
	return nil
}

// configureRouter Объявляем список доступных роутов
func (s *server) configureRouter() {
	s.router.Use(middleware.RequestID())
	api := s.router.Group("/api", s.ErrorHandler)
	{
		api.GET("/", s.handleVersion)
		api.GET("/status", s.handleStatus)
		api.GET("/swagger/*", echoSwagger.WrapHandler)
		v1 := api.Group("/v1", s.ErrorHandler)
		{
			v1.Use(middleware.Logger())
			v1.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(2)))
			v1.POST("/token", s.handleToken)
			v1.POST("/login", s.login)
			authGroup := v1.Group("")
			{
				authGroup.Use(echojwt.WithConfig(echojwt.Config{
					ParseTokenFunc: s.ParseTokenFunc,
				}))
				authGroup.GET("/typelist", s.getTypeList)
				authGroup.GET("/slot", s.getSlotList)
			}
		}
	}
}

func (s *server) ParseTokenFunc(c echo.Context, tokenString string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SigningKey), nil
	})

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
