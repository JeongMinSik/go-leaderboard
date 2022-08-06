package leaderboard

import (
	"context"
	"fmt"

	"github.com/JeongMinSik/go-leaderboard/pkg/redisstorage"
)

type LeaderBoard struct {
	redisStorage *redisstorage.RedisStorage
}

type User struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
	Rank  int64   `json:"rank"`
}

func New() *LeaderBoard {
	return &LeaderBoard{
		redisStorage: redisstorage.New(),
	}
}

func (lb *LeaderBoard) UserCount(ctx context.Context) (int64, error) {
	count, err := lb.redisStorage.Count(ctx)
	return count, fmt.Errorf("lb.redisStorage.Count: %w", err)
}

func (lb *LeaderBoard) AddUser(ctx context.Context, name string, score float64) error {
	return fmt.Errorf("lb.redisStorage.Add: %w", lb.redisStorage.Add(ctx, name, score))
}

func (lb *LeaderBoard) GetUser(ctx context.Context, name string) (User, error) {
	rank, score, err := lb.redisStorage.Get(ctx, name)
	if err != nil {
		return User{}, fmt.Errorf("lb.redisStorage.Get: %w", err)
	}
	return User{
		Name:  name,
		Score: score,
		Rank:  rank,
	}, nil
}

func ErrorWithStatusCode(err error, statusCode int) error {
	return Error{
		origin:     err,
		statusCode: statusCode,
	}
}

type Error struct {
	origin     error
	statusCode int
}

func (e Error) Error() string {
	return e.origin.Error()
}

func (e Error) Unwrap() error {
	return e.origin
}

func (e Error) StatusCode() int {
	return e.statusCode
}
