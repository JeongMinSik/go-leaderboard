package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/JeongMinSik/go-leaderboard/pkg/leaderboard"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	leaderboard leaderboard.LeaderBoard
}

func Setup(e *echo.Echo) {
	handler := &Handler{*leaderboard.New()}

	e.GET("/users/count", handler.GetUserCount)
	e.GET("/users", handler.GetUser)
	e.POST("/users", handler.AddUser)
}

type messageData struct {
	Message string `json:"message"`
}

func errorJson(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	statusCode := http.StatusInternalServerError
	if s, ok := err.(interface{ StatusCode() int }); ok {
		statusCode = s.StatusCode()
	}

	return c.JSON(statusCode, messageData{err.Error()})
}

func (h *Handler) GetUserCount(c echo.Context) error {
	ctx := context.Background()
	count, err := h.leaderboard.UserCount(ctx)
	if err != nil {
		return errorJson(c, err)
	}

	type UserCountData struct {
		Count int64 `json:"count"`
	}

	return c.JSON(http.StatusOK, UserCountData{Count: count})
}

func (h *Handler) GetUser(c echo.Context) error {
	ctx := context.Background()
	userName := c.Param("name")
	user, err := h.leaderboard.GetUser(ctx, userName)
	if err != nil {
		return errorJson(c, err)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) AddUser(c echo.Context) error {
	ctx := context.Background()
	userName := c.Param("name")
	score, err := strconv.ParseFloat(c.QueryParam("score"), 10)
	if err != nil {
		return c.JSON(http.StatusBadRequest, messageData{"score is empty or invalid format"})
	}

	if err := h.leaderboard.AddUser(ctx, userName, score); err != nil {
		return errorJson(c, err)
	}

	type UserCountData struct {
		Count int64 `json:"count"`
	}

	user, err := h.leaderboard.GetUser(ctx, userName)
	if err != nil {
		return errorJson(c, err)
	}

	return c.JSON(http.StatusOK, user)
}
