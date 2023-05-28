package main

import (
	"context"
	"hack2023/internal/app/config"
	"hack2023/internal/app/server"
	"hack2023/internal/app/store"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// @title           Hack2023
// @version         1.0
// @description     Документация по задаче #2 команды "Just do it" участника хакатона leaders2023.innoagency.ru
// @description
// @description     В API реализована JWT токен OAuth 2.0 модель авторизации, с коротко-живущим ключом access_token и долгоживущим ключом refresh_token
// @description
// @description     Авторизация через header «Authorization: Bearer some_jwt_token»
// @description
// @description     Представитель бизнеса:
// @description     login: user
// @description     pwd: 123321
// @description
// @description     Представитель КНО:
// @description     login: kno
// @description     pwd: 123321

// @host      hack.torbeno.ru
// @BasePath  /api

// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345".
func main() {
	config := config.Get()

	defer func() {
		if msg := recover(); msg != nil {
			log.Println("Panic: ", msg)
		}
	}()

	//подключение к бд
	store, err := store.New(config)
	if err != nil {
		log.Print(err)
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile(config.FireBaseFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Print(err)
	}
	push, err := app.Messaging(ctx)
	if err != nil {
		log.Print(err)
	}

	srv := server.NewServer(store, push, config)

	if err := srv.Start(config); err != nil && err != http.ErrServerClosed {
		log.Print(err)
	}

}
