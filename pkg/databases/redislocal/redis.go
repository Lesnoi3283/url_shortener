package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis() Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var red Redis = Redis{}
	red.Client = client

	return red
}

//
//Save(key string, val string) error
//Get(key string) (string, error)

func (r *Redis) Save(key string, val string) error {
	status := r.Client.Set(context.Background(), key, val, time.Hour*24*31)
	return status.Err()
}

func (r *Redis) Get(key string) (string, error) {
	str, err := r.Client.Get(context.Background(), key)
	return str, err
}
