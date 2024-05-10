package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sknutsen/planner/handler"
	"github.com/sknutsen/planner/router"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Could not load .env file. Err: %s\n", err)
	}

	h := &handler.Handler{
		Port: os.Getenv("PORT"),
	}

	e := echo.New()

	e.Use(middleware.Logger())

	router.Setup(e, h)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", h.Port)))
}
