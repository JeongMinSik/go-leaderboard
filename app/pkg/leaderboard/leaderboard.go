package leaderboard

import (
	"context"
	"net/http"

	"github.com/JeongMinSik/go-leaderboard/pkg/redisstorage"
	"github.com/pkg/errors"
)

type Interface interface {
	UserCount(ctx context.Context) (int64, error)
	AddUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, name string) (*UserRank, error)
	DeleteUser(ctx context.Context, name string) (bool, error)
	UpdateUser(ctx context.Context, user User) error
	GetUserList(ctx context.Context, start int64, stop int64) ([]User, error)
}

type LeaderBoard struct {
	redisStorage *redisstorage.RedisStorage
}

type User struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type UserRank struct {
	User
	Rank int64 `json:"rank"`
}

func New() (Interface, error) {
	db, err := redisstorage.New()
	if err != nil {
		return nil, errors.Wrap(err, "redisstorage.New()")
	}
	if db == nil {
		return nil, errors.New("redis nil")
	}
	return &LeaderBoard{
		redisStorage: db,
	}, nil
}

func (lb *LeaderBoard) UserCount(ctx context.Context) (int64, error) {
	count, err := lb.redisStorage.Count(ctx)
	return count, errors.Wrap(err, "lb.redisStorage.Count")
}

func (lb *LeaderBoard) AddUser(ctx context.Context, user User) error {
	return errors.Wrap(lb.redisStorage.Add(ctx, user.Name, user.Score), "lb.redisStorage.Add")
}

func (lb *LeaderBoard) GetUser(ctx context.Context, name string) (*UserRank, error) {
	exists, rank, score, err := lb.redisStorage.Get(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "lb.redisStorage.Get")
	} else if !exists {
		return nil, ErrorWithStatusCode(errors.New("not exists user: "+name), http.StatusNotFound)
	}
	return &UserRank{
		User: User{
			Name:  name,
			Score: score,
		},
		Rank: rank,
	}, nil
}

func (lb *LeaderBoard) DeleteUser(ctx context.Context, name string) (bool, error) {
	ok, err := lb.redisStorage.Delete(ctx, name)
	return ok, errors.Wrap(err, "lb.redisStorage.Delete")
}

func (lb *LeaderBoard) UpdateUser(ctx context.Context, user User) error {
	exists, err := lb.redisStorage.Update(ctx, user.Name, user.Score)
	if err != nil {
		return errors.Wrap(err, "lb.redisStorage.Update")
	}
	if !exists {
		return ErrorWithStatusCode(errors.New("not exists user:"+user.Name), http.StatusNotFound)
	}
	return nil
}

func (lb *LeaderBoard) GetUserList(ctx context.Context, start int64, stop int64) ([]User, error) {
	userList, err := lb.redisStorage.Range(ctx, start, stop)
	if err != nil {
		return nil, errors.Wrap(err, "lb.redisStorage.Range")
	}
	result := make([]User, 0, len(userList))
	for _, user := range userList {
		result = append(result, User{
			Name:  user.Member.(string),
			Score: user.Score,
		})
	}
	return result, nil
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
