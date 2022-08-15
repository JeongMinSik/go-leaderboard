package main

import (
	_ "github.com/JeongMinSik/go-leaderboard/docs"

	"github.com/JeongMinSik/go-leaderboard/pkg/handler"
	"github.com/JeongMinSik/go-leaderboard/pkg/leaderboard"
	"github.com/JeongMinSik/go-leaderboard/pkg/logger"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	SetupLogger(e)
	lb, err := leaderboard.New()
	if err != nil {
		e.Logger.Fatal(err)
	}
	SetupHandler(e, lb)
	e.Logger.Fatal(e.Start(":6025"))
}

func SetupLogger(e *echo.Echo) {
	log := logger.New()
	if err := log.AddElasticHook(e, "api-log"); err != nil {
		log.Panic(err)
	}
}

func SetupHandler(e *echo.Echo, lb leaderboard.Interface) {
	hdler := handler.Handler{
		Leaderboard: lb,
	}
	e.GET("/", hdler.Hello)
	e.GET("/teapot", hdler.Teapot)
	e.GET("/users/count", hdler.GetUserCount)
	e.GET("/users", hdler.GetUser)
	e.POST("/users", hdler.AddUser)
	e.DELETE("/users", hdler.DeleteUser)
	e.PATCH("/users", hdler.UpdateUser)
	e.GET("/users/:start/to/:stop", hdler.GetUserList)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
