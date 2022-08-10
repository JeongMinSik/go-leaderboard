package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/JeongMinSik/go-leaderboard/pkg/leaderboard"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handler struct {
	leaderboard leaderboard.LeaderBoard
}

func Setup(e *echo.Echo) {
	handler := &Handler{*leaderboard.New()}

	e.GET("/", handler.Hello)
	e.GET("/teapot", handler.Teapot)
	e.GET("/users/count", handler.GetUserCount)
	e.GET("/users", handler.GetUser)
	e.POST("/users", handler.AddUser)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

type messageData struct {
	Message string `json:"message"`
}

type userCountData struct {
	Count int64 `json:"count"`
}

func responseJSON(c echo.Context, statusCode int, data interface{}) error {
	return errors.Wrap(c.JSON(statusCode, data), "c.JSON")
}

func errorJSON(c echo.Context, err error) error {
	if err == nil {
		return nil
	}
	statusCode := http.StatusInternalServerError
	var apiErr interface{ StatusCode() int }
	if ok := errors.As(err, &apiErr); ok {
		statusCode = apiErr.StatusCode()
	}
	return responseJSON(c, statusCode, messageData{err.Error()})
}

// @Description 테스트용
// @Tags        test
// @Success     200 {string} string "Hello go-leaderboard"
// @Router      / [get]
func (h *Handler) Hello(c echo.Context) error {
	return responseJSON(c, http.StatusOK, "Hello go-leaderboard")
}

// @Description 테스트용
// @Tags        test
// @failure     418 {string} string "I'm a teapot"
// @Router      /teapot [get]
func (h *Handler) Teapot(c echo.Context) error {
	return responseJSON(c, http.StatusTeapot, "I'm a teapot")
}

// @Summary      Get user count
// @Description  전체 유저 수
// @Tags         Users
// @Produce      json
// @Success      200  {object}  userCountData
// @Failure      400  {object}  messageData "name 또는 score param 확인 필요"
// @Failure      500  {object}  messageData "서버에러"
// @Router       /users/count [get]
func (h *Handler) GetUserCount(c echo.Context) error {
	ctx := context.Background()
	count, err := h.leaderboard.UserCount(ctx)
	if err != nil {
		return errorJSON(c, err)
	}
	return responseJSON(c, http.StatusOK, userCountData{Count: count})
}

// @Summary      Show a user info
// @Description  name으로 User의 score와 rank를 얻습니다.
// @Tags         Users
// @Produce      json
// @Param        name   query   string  true  "User name"
// @Success      200  {object}  leaderboard.User
// @Failure      400  {object}  messageData "name param 확인 필요"
// @Failure      500  {object}  messageData "서버에러"
// @Router       /users [get]
func (h *Handler) GetUser(c echo.Context) error {
	ctx := context.Background()
	userName := c.QueryParam("name")
	if userName == "" {
		return responseJSON(c, http.StatusBadRequest, messageData{"user name is empty"})
	}
	user, err := h.leaderboard.GetUser(ctx, userName)
	if err != nil {
		return errorJSON(c, err)
	}
	return responseJSON(c, http.StatusOK, user)
}

// @Summary      Add a user
// @Description  신규 user를 추가합니다.
// @Tags         Users
// @Produce      json
// @Param        name   query   string  true  "User name"
// @Param        score  query   number  true  "User score"
// @Success      200  {object}  leaderboard.User
// @Failure      400  {object}  messageData "name 또는 score param 확인 필요"
// @Failure      500  {object}  messageData "서버에러"
// @Router       /users [post]
func (h *Handler) AddUser(c echo.Context) error {
	ctx := context.Background()
	userName := c.QueryParam("name")
	if userName == "" {
		return responseJSON(c, http.StatusBadRequest, messageData{"user name is empty"})
	}
	score, err := strconv.ParseFloat(c.QueryParam("score"), 64)
	if err != nil {
		return responseJSON(c, http.StatusBadRequest, messageData{"score is empty or invalid format"})
	}
	if err := h.leaderboard.AddUser(ctx, userName, score); err != nil {
		return errorJSON(c, err)
	}
	user, err := h.leaderboard.GetUser(ctx, userName)
	if err != nil {
		return errorJSON(c, err)
	}
	return responseJSON(c, http.StatusOK, user)
}
