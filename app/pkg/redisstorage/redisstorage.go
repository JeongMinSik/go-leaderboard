package redisstorage

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisStorage struct {
	zsetKey string
	client  *redis.Client
}

func New() (*RedisStorage, error) {
	zsetKey := "scores"
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, errors.New("empty redis addr")
	}
	db := redis.NewClient(&redis.Options{Addr: redisAddr})
	if db == nil {
		return nil, errors.New("redis client is nil")
	}

	return &RedisStorage{
		zsetKey: zsetKey,
		client:  db,
	}, nil
}

func NewMock(zsetKey string, db *redis.Client) *RedisStorage {
	return &RedisStorage{
		zsetKey: zsetKey,
		client:  db,
	}
}

func (r *RedisStorage) Add(ctx context.Context, name string, score float64) error {
	addCount, err := r.client.ZAddNX(ctx, r.zsetKey, &redis.Z{Score: score, Member: name}).Result()
	if err != nil {
		return errors.Wrap(err, "r.client.ZAddNX")
	}
	if addCount != 1 {
		return errors.New("already exists user:" + name)
	}
	return nil
}

func (r *RedisStorage) Count(ctx context.Context) (int64, error) {
	count, err := r.client.ZCount(ctx, r.zsetKey, "-inf", "+inf").Result()
	return count, errors.Wrap(err, "ZCount")
}

func (r *RedisStorage) Get(ctx context.Context, name string) (bool, int64, float64, error) {
	pipe := r.client.TxPipeline()
	scoreCmd := pipe.ZScore(ctx, r.zsetKey, name)
	rankCmd := pipe.ZRank(ctx, r.zsetKey, name)
	if _, err := pipe.Exec(ctx); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, -1, 0.0, nil
		}
		return false, -1, 0.0, errors.Wrap(err, "pipe.Exec")
	}
	score, err := scoreCmd.Result()
	if err != nil {
		return false, -1, 0.0, errors.Wrap(err, "scoreCmd.Result")
	}

	rank, err := rankCmd.Result()
	if err != nil {
		return false, -1, 0.0, errors.Wrap(err, "rankCmd.Result")
	}

	return true, rank, score, nil
}

func (r *RedisStorage) Delete(ctx context.Context, name string) (bool, error) {
	remCount, err := r.client.ZRem(ctx, r.zsetKey, name).Result()
	return remCount == 1, errors.Wrap(err, "ZRem")
}

func (r *RedisStorage) Update(ctx context.Context, name string, score float64) (bool, error) {
	updateCount, err := r.client.ZAddXX(ctx, r.zsetKey, &redis.Z{Score: score, Member: name}).Result()
	if err != nil {
		return false, errors.Wrap(err, "r.client.ZAddXX")
	}
	return updateCount == 1, nil
}

func (r *RedisStorage) Range(ctx context.Context, start int64, stop int64) ([]redis.Z, error) {
	userList, err := r.client.ZRangeWithScores(ctx, r.zsetKey, start, stop).Result()
	if err != nil {
		return nil, errors.Wrap(err, "r.client.ZRange")
	}

	return userList, nil
}
