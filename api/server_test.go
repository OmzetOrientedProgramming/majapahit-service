package api

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service.git/internal/checkup"
)

func TestInitServer(t *testing.T) {
	_ = godotenv.Load("../.env")
	e := echo.New()
	server := NewServer(e)

	server.Init()

	assert.NotEqual(t, checkup.NewRepo(nil), checkUpRepo)
	assert.NotEqual(t, checkup.NewService(nil), checkUpService)
	assert.NotEqual(t, checkup.NewHandler(nil), checkupHandler)
}

func TestRunServerFailed(t *testing.T) {
	e := echo.New()
	server := NewServer(e)

	server.RunServer("testport")
}

func TestRunServerSuccess(t *testing.T) {
	e := echo.New()
	server := NewServer(e)
	done := make(chan bool)

	go func() {
		server.RunServer("8080")
		done <- true
	}()
	server.Router.Shutdown(context.Background())
}
