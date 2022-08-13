package main

import (
	_ "github.com/JeongMinSik/go-leaderboard/docs"

	"github.com/JeongMinSik/go-leaderboard/pkg/handler"
	"github.com/JeongMinSik/go-leaderboard/pkg/leaderboard"
	"github.com/labstack/echo/v4"
)

// @title       Leaderboard API
// @version     1.0
// @description go언어로 만든 리더보드 api 토이프로젝트입니다.

// @contact.name  JeongMinSik
// @contact.url   https://github.com/JeongMinSik/go-leaderboard
// @contact.email jms6025a@naver.com

// @host localhost:6025
func main() {
	e := echo.New()
	lb, err := leaderboard.New()
	if err != nil {
		e.Logger.Fatal(err)
	}
	hd := handler.Handler{
		Leaderboard: lb,
	}
	handler.Setup(e, hd)
	e.Logger.Fatal(e.Start(":6025"))
}
