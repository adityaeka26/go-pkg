package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Db *redis.Client
}

func NewRedis(host, password string) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	err := rdb.Set(ctx, "ping", "pong", time.Second*10).Err()
	if err != nil {
		return nil, err
	}
	_, err = rdb.Get(ctx, "ping").Result()
	if err != nil {
		return nil, err
	}
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &Redis{
		Db: rdb,
	}, nil
}

func (r *Redis) Close(ctx context.Context) error {
	return r.Db.Close()
}

func (r *Redis) Set(key, value string, expiration time.Duration, ctx context.Context) error {
	err := r.Db.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Get(key string, ctx context.Context) (*string, error) {
	val, err := r.Db.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return &val, nil
}

func (r *Redis) TTL(key string, ctx context.Context) (time.Duration, error) {
	ttl, err := r.Db.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return ttl, nil
}

func (r *Redis) Del(key string, ctx context.Context) error {
	err := r.Db.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
