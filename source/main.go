package main

import (
	"github.com/JeongMinSik/go-leaderboard/pkg/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	handler.Setup(e)
	e.Logger.Fatal(e.Start(":6025"))
}
