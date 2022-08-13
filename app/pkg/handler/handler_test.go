package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JeongMinSik/go-leaderboard/pkg/leaderboard"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wangjia184/sortedset"
)

type TestLeaderBoard struct {
	userSet sortedset.SortedSet
}

func (lb *TestLeaderBoard) UserCount(_ context.Context) (int64, error) {
	return int64(lb.userSet.GetCount()), nil
}

func (lb *TestLeaderBoard) AddUser(_ context.Context, user leaderboard.User) error {
	if data := lb.userSet.GetByKey(user.Name); data != nil {
		return errors.New("already exists name: " + user.Name)
	}
	if ok := lb.userSet.AddOrUpdate(user.Name, sortedset.SCORE(user.Score), nil); !ok {
		return errors.New("update data: " + user.Name)
	}
	return nil
}

func (lb *TestLeaderBoard) GetUser(_ context.Context, name string) (*leaderboard.UserRank, error) {
	rank := lb.userSet.FindRank(name)
	if rank == 0 {
		return nil, errors.New("not exists name: " + name)
	}
	node := lb.userSet.GetByKey(name)
	if node == nil {
		return nil, errors.New("not exists name: " + name)
	}
	return &leaderboard.UserRank{
		User: leaderboard.User{
			Name:  name,
			Score: float64(node.Score()),
		},
		// redis zset의 rank는 오름차순으로 0부터 시작하므로 보정
		Rank: int64(lb.userSet.GetCount() - rank),
	}, nil
}

func (lb *TestLeaderBoard) DeleteUser(_ context.Context, name string) (bool, error) {
	return lb.userSet.Remove(name) != nil, nil
}

func (lb *TestLeaderBoard) UpdateUser(_ context.Context, user leaderboard.User) error {
	node := lb.userSet.GetByKey(user.Name)
	if node == nil {
		return errors.New("not exists name: " + user.Name)
	}
	if ok := lb.userSet.AddOrUpdate(user.Name, sortedset.SCORE(user.Score), nil); ok {
		return errors.New("new name: " + user.Name)
	}
	return nil
}

func (lb *TestLeaderBoard) GetUserList(_ context.Context, start int64, stop int64) ([]leaderboard.User, error) {
	result := make([]leaderboard.User, 0, lb.userSet.GetCount())
	nodes := lb.userSet.GetByRankRange(int(start), int(stop), false)
	// sortedSet은 오름차순, redis zset은 내림차순이므로 역순으로 반환
	for i := len(nodes) - 1; i >= 0; i-- {
		result = append(result, leaderboard.User{
			Name:  nodes[i].Key(),
			Score: float64(nodes[i].Score()),
		})
	}
	return result, nil
}

func TestAddUser(t *testing.T) {
	// Setup
	e := echo.New()
	h := &Handler{&TestLeaderBoard{
		userSet: *sortedset.New(),
	}}

	// AddUser
	const userJSON = `{"name": "Minsik", "score": 100, "rank": 0}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, h.AddUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		require.JSONEq(t, userJSON, rec.Body.String())
	}
}

func TestGetUser(t *testing.T) {
	// Setup
	e := echo.New()
	sortedSet := sortedset.New()
	sortedSet.AddOrUpdate("Minsik", 10000, nil)
	h := &Handler{&TestLeaderBoard{
		userSet: *sortedSet,
	}}

	// GetUser
	req := httptest.NewRequest(http.MethodGet, "/users?name=Minsik", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, h.GetUser(c)) {
		const userJSON = `{"name": "Minsik", "score": 10000, "rank":0}`
		assert.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, userJSON, rec.Body.String())
	}
}

func TestDeleteUser(t *testing.T) {
	// Setup
	e := echo.New()
	sortedSet := sortedset.New()
	sortedSet.AddOrUpdate("Minsik", 10000, nil)
	h := &Handler{&TestLeaderBoard{
		userSet: *sortedSet,
	}}

	// GetUserCount
	req := httptest.NewRequest(http.MethodGet, "/users?name=Minsik", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, h.GetUserCount(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		const countJSON = `{"count": 1}`
		require.JSONEq(t, countJSON, rec.Body.String())
	}

	// DeleteUser
	req2 := httptest.NewRequest(http.MethodDelete, "/users?name=Minsik", nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	if assert.NoError(t, h.DeleteUser(c2)) {
		assert.Equal(t, http.StatusOK, rec2.Code)
		const deleteJSON = `{"name": "Minsik", "is_deleted": true}`
		require.JSONEq(t, deleteJSON, rec2.Body.String())
	}

	// GetUserCount
	req3 := httptest.NewRequest(http.MethodGet, "/users?name=Minsik", nil)
	rec3 := httptest.NewRecorder()
	c3 := e.NewContext(req3, rec3)
	if assert.NoError(t, h.GetUserCount(c3)) {
		assert.Equal(t, http.StatusOK, rec3.Code)
		const countJSON = `{"count": 0}`
		require.JSONEq(t, countJSON, rec3.Body.String())
	}
}

func TestUpdateUser(t *testing.T) {
	// Setup
	e := echo.New()
	sortedSet := sortedset.New()
	sortedSet.AddOrUpdate("Yumi", 500, nil)
	sortedSet.AddOrUpdate("Minsik", 100, nil)
	h := &Handler{&TestLeaderBoard{
		userSet: *sortedSet,
	}}

	// GetUser
	req := httptest.NewRequest(http.MethodGet, "/users?name=Minsik", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, h.GetUser(c)) {
		const userJSON = `{"name": "Minsik", "score": 100, "rank":1}`
		assert.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, userJSON, rec.Body.String())
	}

	// UpdateUser
	const reqUserJSON = `{"name": "Minsik", "score": 10000}`
	req2 := httptest.NewRequest(http.MethodPatch, "/users", strings.NewReader(reqUserJSON))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	if assert.NoError(t, h.UpdateUser(c2)) {
		const newUserJSON = `{"name": "Minsik", "score": 10000, "rank":0}`
		assert.Equal(t, http.StatusOK, rec2.Code)
		require.JSONEq(t, newUserJSON, rec2.Body.String())
	}
}
