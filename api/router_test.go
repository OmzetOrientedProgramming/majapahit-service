package api

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRouter(t *testing.T) {
	e := echo.New()
	server := NewServer(e)
	server.Init()
}
