package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JeongMinSik/go-leaderboard/pkg/leaderboard"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userJSON = `{"name": "Minsik", "score": 10000, "rank": 0}`
)

type TestLeaderBoard struct {
}

func (lb *TestLeaderBoard) UserCount(_ context.Context) (int64, error) {
	return 0, nil
}

func (lb *TestLeaderBoard) AddUser(_ context.Context, _ leaderboard.User) error {
	return nil
}

func (lb *TestLeaderBoard) GetUser(_ context.Context, name string) (leaderboard.UserRank, error) {
	return leaderboard.UserRank{
		User: leaderboard.User{
			Name:  name,
			Score: 10000,
		},
		Rank: 0,
	}, nil
}

func (lb *TestLeaderBoard) DeleteUser(_ context.Context, _ string) (bool, error) {
	return true, nil
}

func (lb *TestLeaderBoard) UpdateUser(_ context.Context, _ leaderboard.User) error {
	return nil
}

func (lb *TestLeaderBoard) GetUserList(_ context.Context, _ int64, _ int64) ([]leaderboard.User, error) {
	return nil, nil
}

func TestAddUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &Handler{&TestLeaderBoard{}}

	// Assertions
	if assert.NoError(t, h.AddUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		require.JSONEq(t, userJSON, rec.Body.String())
	}
}

func TestGetUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	q := req.URL.Query()
	q.Add("name", "Minsik")
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users")
	c.SetParamNames()
	h := &Handler{&TestLeaderBoard{}}

	// Assertions
	if assert.NoError(t, h.GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, userJSON, rec.Body.String())
	}
}
