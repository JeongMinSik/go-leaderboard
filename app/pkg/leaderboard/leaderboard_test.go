package leaderboard

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/JeongMinSik/go-leaderboard/pkg/redisstorage"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

const ZSetKeyName = "scores"

func TestNew(t *testing.T) {
	_, err := New()
	assert.ErrorContains(t, err, "empty redis addr")
}

func TestErrorWithStatusCode(t *testing.T) {
	origin := errors.New("test error")
	err := ErrorWithStatusCode(origin, http.StatusBadRequest)
	var apiErr interface {
		StatusCode() int
		Error() string
		Unwrap() error
	}
	if assert.ErrorAs(t, err, &apiErr) {
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
		assert.Equal(t, "test error", apiErr.Error())
		assert.Equal(t, origin, apiErr.Unwrap())
	}
}

func TestUserCount(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	lb := &LeaderBoard{
		redisStorage: redisstorage.NewMock(ZSetKeyName, db),
	}

	mock.ExpectZCount(ZSetKeyName, "-inf", "+inf").SetVal(3)

	userCount, err := lb.UserCount(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, int64(3), userCount)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	lb := &LeaderBoard{
		redisStorage: redisstorage.NewMock(ZSetKeyName, db),
	}

	mock.ExpectZAddNX(ZSetKeyName, &redis.Z{
		Score:  100,
		Member: "Minsik",
	}).SetVal(1)

	err := lb.AddUser(ctx, User{
		Name:  "Minsik",
		Score: 100,
	})
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	lb := &LeaderBoard{
		redisStorage: redisstorage.NewMock(ZSetKeyName, db),
	}

	mock.ExpectTxPipeline()
	mock.ExpectZScore(ZSetKeyName, "Minsik").SetVal(999)
	mock.ExpectZRank(ZSetKeyName, "Minsik").SetVal(4)
	mock.ExpectTxPipelineExec()

	userRank, err := lb.GetUser(ctx, "Minsik")
	if assert.NoError(t, err) {
		assert.Equal(t, UserRank{
			User: User{
				Name:  "Minsik",
				Score: 999,
			},
			Rank: 4,
		}, *userRank)
	}

	mock.ExpectTxPipeline()
	mock.ExpectZScore(ZSetKeyName, "Foo").RedisNil()
	_, err = lb.GetUser(ctx, "Foo")
	var apiErr interface{ StatusCode() int }
	if assert.ErrorAs(t, err, &apiErr) {
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	lb := &LeaderBoard{
		redisStorage: redisstorage.NewMock(ZSetKeyName, db),
	}

	mock.ExpectZRem(ZSetKeyName, "Minsik").SetVal(1)

	ok, err := lb.DeleteUser(ctx, "Minsik")

	if assert.NoError(t, err) {
		assert.Equal(t, true, ok)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	lb := &LeaderBoard{
		redisStorage: redisstorage.NewMock(ZSetKeyName, db),
	}

	mock.ExpectZAddXX(ZSetKeyName, &redis.Z{
		Score:  100,
		Member: "Minsik",
	}).SetVal(1)

	err := lb.UpdateUser(ctx, User{
		Name:  "Minsik",
		Score: 100,
	})

	assert.NoError(t, err)

	mock.ExpectZAddXX(ZSetKeyName, &redis.Z{
		Score:  200,
		Member: "Foo",
	}).SetVal(0)

	err = lb.UpdateUser(ctx, User{
		Name:  "Foo",
		Score: 200,
	})

	var apiErr interface{ StatusCode() int }
	if assert.True(t, errors.As(err, &apiErr)) {
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestUserList(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	lb := &LeaderBoard{
		redisStorage: redisstorage.NewMock(ZSetKeyName, db),
	}

	mock.ExpectZRangeWithScores(ZSetKeyName, 0, 2).SetVal([]redis.Z{
		{
			Score:  1000,
			Member: "Minsik",
		},
		{
			Score:  500,
			Member: "Foo",
		},
		{
			Score:  100,
			Member: "FooFoo",
		},
	})

	users, err := lb.GetUserList(ctx, 0, 2)

	if assert.NoError(t, err) {
		expected := []User{
			{
				Name:  "Minsik",
				Score: 1000,
			},
			{
				Name:  "Foo",
				Score: 500,
			},
			{
				Name:  "FooFoo",
				Score: 100,
			},
		}
		for i, user := range users {
			assert.Equal(t, expected[i], user)
		}
	}

	mock.ExpectZRangeWithScores(ZSetKeyName, 0, 1).SetErr(redis.ErrClosed)
	_, err = lb.GetUserList(ctx, 0, 1)
	assert.Error(t, err, redis.ErrClosed)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
