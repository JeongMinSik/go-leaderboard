package leaderboard

import (
	"context"

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
	return lb.redisStorage.Count(ctx)
}

func (lb *LeaderBoard) AddUser(ctx context.Context, name string, score float64) error {
	return lb.redisStorage.Add(ctx, name, score)
}

func (lb *LeaderBoard) GetUser(ctx context.Context, name string) (User, error) {
	rank, score, err := lb.redisStorage.Get(ctx, name)
	if err != nil {
		return User{}, err
	}
	return User{
		Name:  name,
		Score: score,
		Rank:  rank,
	}, nil
}
