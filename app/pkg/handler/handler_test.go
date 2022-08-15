package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JeongMinSik/go-leaderboard/pkg/leaderboard"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wangjia184/sortedset"
)

type FakeLeaderBoard struct {
	UserSet sortedset.SortedSet
}

func (lb *FakeLeaderBoard) UserCount(_ context.Context) (int64, error) {
	return int64(lb.UserSet.GetCount()), nil
}

func (lb *FakeLeaderBoard) AddUser(_ context.Context, user leaderboard.User) error {
	if data := lb.UserSet.GetByKey(user.Name); data != nil {
		err := leaderboard.ErrorWithStatusCode(errors.New("already exists name: "+user.Name), http.StatusBadRequest)
		return errors.Wrap(err, "lb.UserSet.GetByKey")
	}
	if ok := lb.UserSet.AddOrUpdate(user.Name, sortedset.SCORE(user.Score), nil); !ok {
		err := leaderboard.ErrorWithStatusCode(errors.New("update data: "+user.Name), http.StatusInternalServerError)
		return errors.Wrap(err, "lb.UserSet.AddOrUpdate")
	}
	return nil
}

func (lb *FakeLeaderBoard) GetUser(_ context.Context, name string) (*leaderboard.UserRank, error) {
	rank := lb.UserSet.FindRank(name)
	if rank == 0 {
		err := leaderboard.ErrorWithStatusCode(errors.New("not exists name: "+name), http.StatusBadRequest)
		return nil, errors.Wrap(err, "lb.UserSet.FindRank")
	}
	node := lb.UserSet.GetByKey(name)
	if node == nil {
		err := leaderboard.ErrorWithStatusCode(errors.New("not exists name: "+name), http.StatusBadRequest)
		return nil, errors.Wrap(err, "lb.UserSet.GetByKey")
	}
	return &leaderboard.UserRank{
		User: leaderboard.User{
			Name:  name,
			Score: float64(node.Score()),
		},
		// redis zset의 rank는 오름차순으로 0부터 시작하므로 보정
		Rank: int64(lb.UserSet.GetCount() - rank),
	}, nil
}

func (lb *FakeLeaderBoard) DeleteUser(_ context.Context, name string) (bool, error) {
	return lb.UserSet.Remove(name) != nil, nil
}

func (lb *FakeLeaderBoard) UpdateUser(_ context.Context, user leaderboard.User) error {
	node := lb.UserSet.GetByKey(user.Name)
	if node == nil {
		err := leaderboard.ErrorWithStatusCode(errors.New("not exists user: "+user.Name), http.StatusBadRequest)
		return errors.Wrap(err, "lb.UserSet.GetByKey")
	}
	if ok := lb.UserSet.AddOrUpdate(user.Name, sortedset.SCORE(user.Score), nil); ok {
		err := leaderboard.ErrorWithStatusCode(errors.New("new name: "+user.Name), http.StatusInternalServerError)
		return errors.Wrap(err, "lb.UserSet.AddOrUpdate")
	}
	return nil
}

func (lb *FakeLeaderBoard) GetUserList(_ context.Context, start int64, stop int64) ([]leaderboard.User, error) {
	result := make([]leaderboard.User, 0, lb.UserSet.GetCount())
	nodes := lb.UserSet.GetByRankRange(int(start+1), int(stop+1), false) // index + 1 보정
	// sortedSet은 오름차순, redis zset은 내림차순이므로 역순으로 반환
	for i := len(nodes) - 1; i >= 0; i-- {
		result = append(result, leaderboard.User{
			Name:  nodes[i].Key(),
			Score: float64(nodes[i].Score()),
		})
	}
	return result, nil
}

func TestErrorJSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	err := leaderboard.ErrorWithStatusCode(errors.New("test error"), http.StatusBadRequest)
	if assert.NoError(t, errorJSON(ctx, err)) {
		const errorJSON = `{"message": "test error"}`
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		require.JSONEq(t, errorJSON, rec.Body.String())
	}
	assert.NoError(t, errorJSON(ctx, nil))
}

func TestHello(t *testing.T) {
	// Setup
	e := echo.New()
	h := &Handler{}

	// AddUser
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, h.Hello(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "\"Hello go-leaderboard\"\n", rec.Body.String())
	}
}

func TestTeapot(t *testing.T) {
	// Setup
	e := echo.New()
	h := &Handler{}

	// AddUser
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, h.Teapot(c)) {
		assert.Equal(t, http.StatusTeapot, rec.Code)
		assert.Equal(t, "\"I'm a teapot\"\n", rec.Body.String())
	}
}

func TestAddUser(t *testing.T) {
	// Setup
	e := echo.New()
	h := &Handler{&FakeLeaderBoard{
		UserSet: *sortedset.New(),
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

	// AddUser - error
	const invalidJSON = `{"name": "Mins`
	req2 := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(invalidJSON))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	if assert.NoError(t, h.AddUser(c2)) {
		const errorJSON = `{"message": "invalid body: user info"}`
		assert.Equal(t, http.StatusBadRequest, rec2.Code)
		require.JSONEq(t, errorJSON, rec2.Body.String())
	}
}

func TestGetUser(t *testing.T) {
	// Setup
	e := echo.New()
	sortedSet := sortedset.New()
	sortedSet.AddOrUpdate("Minsik", 10000, nil)
	h := &Handler{&FakeLeaderBoard{
		UserSet: *sortedSet,
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

	// GetUser - not exists
	req2 := httptest.NewRequest(http.MethodGet, "/users?name=Foo", nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	if assert.NoError(t, h.GetUser(c2)) {
		const errorJSON = `{"message": "lb.UserSet.FindRank: not exists name: Foo"}`
		assert.Equal(t, http.StatusBadRequest, rec2.Code)
		require.JSONEq(t, errorJSON, rec2.Body.String())
	}

	// GetUser - empty name
	req3 := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec3 := httptest.NewRecorder()
	c3 := e.NewContext(req3, rec3)
	if assert.NoError(t, h.GetUser(c3)) {
		const errorJSON = `{"message": "user name is empty"}`
		assert.Equal(t, http.StatusBadRequest, rec3.Code)
		require.JSONEq(t, errorJSON, rec3.Body.String())
	}
}

func TestDeleteUser(t *testing.T) {
	// Setup
	e := echo.New()
	sortedSet := sortedset.New()
	sortedSet.AddOrUpdate("Minsik", 10000, nil)
	h := &Handler{&FakeLeaderBoard{
		UserSet: *sortedSet,
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

	// DeleteUser - error
	req4 := httptest.NewRequest(http.MethodDelete, "/users", nil)
	rec4 := httptest.NewRecorder()
	c4 := e.NewContext(req4, rec4)
	if assert.NoError(t, h.DeleteUser(c4)) {
		const errorJSON = `{"message": "user name is empty"}`
		assert.Equal(t, http.StatusBadRequest, rec4.Code)
		require.JSONEq(t, errorJSON, rec4.Body.String())
	}
}

func TestUpdateUser(t *testing.T) {
	// Setup
	e := echo.New()
	sortedSet := sortedset.New()
	sortedSet.AddOrUpdate("Yumi", 500, nil)
	sortedSet.AddOrUpdate("Minsik", 100, nil)
	h := &Handler{&FakeLeaderBoard{
		UserSet: *sortedSet,
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

	// UpdateUser - invalid body
	const invalidJSON = `{"name": "Minsik", "scor`
	req3 := httptest.NewRequest(http.MethodPatch, "/users", strings.NewReader(invalidJSON))
	req3.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec3 := httptest.NewRecorder()
	c3 := e.NewContext(req3, rec3)
	if assert.NoError(t, h.UpdateUser(c3)) {
		const errorJSON = `{"message": "invalid body: user info"}`
		assert.Equal(t, http.StatusBadRequest, rec3.Code)
		require.JSONEq(t, errorJSON, rec3.Body.String())
	}

	// UpdateUser - not exists
	const notExistsUserJSON = `{"name": "FooFoo", "score": 200}`
	req4 := httptest.NewRequest(http.MethodPatch, "/users", strings.NewReader(notExistsUserJSON))
	req4.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec4 := httptest.NewRecorder()
	c4 := e.NewContext(req4, rec4)
	if assert.NoError(t, h.UpdateUser(c4)) {
		const errorJSON = `{"message": "lb.UserSet.GetByKey: not exists user: FooFoo"}`
		assert.Equal(t, http.StatusBadRequest, rec4.Code)
		require.JSONEq(t, errorJSON, rec4.Body.String())
	}
}

func TestUserList(t *testing.T) {
	// Setup
	e := echo.New()
	sortedSet := sortedset.New()
	sortedSet.AddOrUpdate("Yumi", 500, nil)
	sortedSet.AddOrUpdate("Minsik", 100, nil)
	sortedSet.AddOrUpdate("Foo", 200, nil)
	sortedSet.AddOrUpdate("FooFoo", 300, nil)
	h := &Handler{&FakeLeaderBoard{
		UserSet: *sortedSet,
	}}

	// GetUserList 1
	req := httptest.NewRequest(http.MethodGet, "/users/:start/to/:stop", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("start", "stop")
	c.SetParamValues("1", "2")
	if assert.NoError(t, h.GetUserList(c)) {
		const userListJSON = `[
			{"name": "FooFoo", "score": 300},
			{"name": "Foo", "score": 200}
		]`
		assert.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, userListJSON, rec.Body.String())
	}

	// GetUserList 2
	req2 := httptest.NewRequest(http.MethodGet, "/users/:start/to/:stop", nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	c2.SetParamNames("start", "stop")
	c2.SetParamValues("0", "3")
	if assert.NoError(t, h.GetUserList(c2)) {
		const userListJSON = `[
			{"name": "Yumi", "score": 500},
			{"name": "FooFoo", "score": 300},
			{"name": "Foo", "score": 200},
			{"name": "Minsik", "score": 100}
		]`
		assert.Equal(t, http.StatusOK, rec2.Code)
		require.JSONEq(t, userListJSON, rec2.Body.String())
	}

	// GetUserList - invalid start
	req3 := httptest.NewRequest(http.MethodGet, "/users/:start/to/:stop", nil)
	rec3 := httptest.NewRecorder()
	c3 := e.NewContext(req3, rec3)
	c3.SetParamNames("start", "stop")
	c3.SetParamValues("abc", "3")
	if assert.NoError(t, h.GetUserList(c3)) {
		const errorJSON = `{"message": "invalid start index"}`
		assert.Equal(t, http.StatusBadRequest, rec3.Code)
		require.JSONEq(t, errorJSON, rec3.Body.String())
	}

	// GetUserList - invalid end
	req4 := httptest.NewRequest(http.MethodGet, "/users/:start/to/:stop", nil)
	rec4 := httptest.NewRecorder()
	c4 := e.NewContext(req4, rec4)
	c4.SetParamNames("start", "stop")
	c4.SetParamValues("2", "xyz")
	if assert.NoError(t, h.GetUserList(c4)) {
		const errorJSON = `{"message": "invalid stop index"}`
		assert.Equal(t, http.StatusBadRequest, rec4.Code)
		require.JSONEq(t, errorJSON, rec4.Body.String())
	}
}
