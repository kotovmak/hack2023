package main

import (
	"hack2023/internal/app/config"
	"hack2023/internal/app/server"
	"hack2023/internal/app/store"
	"log"
	"net/http"
)

// @title           Hack2023
// @version         1.0
// @description     API for flutter app

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

	srv := server.NewServer(store, config)

	if err := srv.Start(config); err != nil && err != http.ErrServerClosed {
		log.Print(err)
	}

}
