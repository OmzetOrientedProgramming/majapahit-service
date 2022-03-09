//go:build !test
// +build !test

package main

import (
	"github.com/labstack/echo/v4/middleware"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/api"
)

func main() {
	// Load env var
	err := godotenv.Load()
	if err != nil {
		logrus.Error(".env not found, will use default env")
	}

	// Creating router
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	s := api.NewServer(router)
	s.Init()

	// Running server
	s.RunServer(os.Getenv("PORT"))
}
