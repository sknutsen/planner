package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/handler"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/router"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Could not load .env file. Err: %s\n", err)
	}

	port := os.Getenv("PORT")

	var callbackUrl string

	host := os.Getenv("HOST")

	if host != "" {
		callbackUrl = fmt.Sprintf("http://%s/%s", os.Getenv("HOST"), os.Getenv("AUTH0_CALLBACK_URL"))
	} else {
		var hostName string

		if port == "" {
			port = "8081"
			hostName = "127.0.0.1"
		} else {
			hostName = "0.0.0.0"
		}

		callbackUrl = fmt.Sprintf("http://%s:%s/%s", hostName, port, os.Getenv("AUTH0_CALLBACK_URL"))
	}

	h := &handler.Handler{
		Host: host,
		Port: port,
		AuthConfig: handler.AuthConfig{
			Domain:       os.Getenv("AUTH0_DOMAIN"),
			Audience:     os.Getenv("AUTH0_AUDIENCE"),
			ClientId:     os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
			CallbackUrl:  callbackUrl,
		},
		TursoConfig: handler.TursoConfig{
			PrimaryUrl: os.Getenv("TURSO_DATABASE_URL"),
			AuthToken:  os.Getenv("TURSO_AUTH_TOKEN"),
		},
	}

	h.Setup()

	gob.Register(map[string]interface{}{})
	gob.Register(models.UserProfile{})

	e := echo.New()

	router.Setup(e, h)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", h.Port)))
}
