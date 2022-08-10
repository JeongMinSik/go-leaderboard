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

func New() *RedisStorage {
	zsetKey := "scores"
	redisAddr := os.Getenv("REDIS_ADDR")

	return &RedisStorage{
		zsetKey: zsetKey,
		client:  redis.NewClient(&redis.Options{Addr: redisAddr}),
	}
}

func (r *RedisStorage) Add(ctx context.Context, name string, score float64) error {
	txf := func(tx *redis.Tx) error {
		exists, err := tx.Exists(ctx, name).Result()
		if err != nil {
			return errors.Wrap(err, "tx.Exists")
		}
		if exists != 0 {
			return errors.New("user already exists")
		}

		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.ZAdd(ctx, r.zsetKey, &redis.Z{Score: score, Member: name})
			return nil
		})
		return errors.Wrap(err, "txf")
	}

	return errors.Wrap(r.client.Watch(ctx, txf, r.zsetKey), "redis watch")
}

func (r *RedisStorage) Count(ctx context.Context) (int64, error) {
	count, err := r.client.ZCount(ctx, r.zsetKey, "-inf", "+inf").Result()
	return count, errors.Wrap(err, "ZCount")
}

func (r *RedisStorage) Get(ctx context.Context, name string) (int64, float64, error) {
	pipe := r.client.TxPipeline()
	scoreCmd := pipe.ZScore(ctx, r.zsetKey, name)
	rankCmd := pipe.ZRank(ctx, r.zsetKey, name)
	if _, err := pipe.Exec(ctx); err != nil {
		return -1, 0.0, errors.Wrap(err, "pipe.Exec")
	}
	score, err := scoreCmd.Result()
	if err != nil {
		return -1, 0.0, errors.Wrap(err, "scoreCmd.Result")
	}

	rank, err := rankCmd.Result()
	if err != nil {
		return -1, 0.0, errors.Wrap(err, "rankCmd.Result")
	}

	return rank, score, nil
}

func (r *RedisStorage) Delete(ctx context.Context, name string) (bool, error) {
	rem_n, err := r.client.ZRem(ctx, r.zsetKey, name).Result()
	return rem_n == 1, errors.Wrap(err, "ZRem")
}
