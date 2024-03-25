package redislocal

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	red := &Redis{Client: client}

	return red
}

func (r *Redis) Save(ctx context.Context, key string, val string) error {
	status := r.Client.Set(ctx, key, val, time.Hour*24*31)
	return status.Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	str, err := r.Client.Get(ctx, key).Result()
	return str, err
}
