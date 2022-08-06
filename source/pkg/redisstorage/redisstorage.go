package redisstorage

import (
	"context"
	"errors"
	"os"

	"github.com/go-redis/redis/v8"
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
			return err
		}
		if exists != 0 {
			return errors.New("user already exists")
		}

		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.ZAdd(ctx, r.zsetKey, &redis.Z{Score: score, Member: name})
			return nil
		})
		return err
	}

	return r.client.Watch(ctx, txf, r.zsetKey)
}

func (r *RedisStorage) Count(ctx context.Context) (int64, error) {
	return r.client.ZCount(ctx, r.zsetKey, "-inf", "+inf").Result()
}

func (r *RedisStorage) Get(ctx context.Context, name string) (int64, float64, error) {
	pipe := r.client.TxPipeline()
	scoreCmd := pipe.ZScore(ctx, r.zsetKey, name)
	rankCmd := pipe.ZRank(ctx, r.zsetKey, name)
	if _, err := pipe.Exec(ctx); err != nil {
		return -1, 0.0, err
	}
	score, err := scoreCmd.Result()
	if err != nil {
		return -1, 0.0, err
	}

	rank, err := rankCmd.Result()
	if err != nil {
		return -1, 0.0, err
	}

	return rank, score, nil
}
