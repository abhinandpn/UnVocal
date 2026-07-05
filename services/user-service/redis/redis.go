package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

var (
	Client *goredis.Client
	Ctx    = context.Background()
)

func Connect() error {
	Client = goredis.NewClient(&goredis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return Client.Ping(Ctx).Err()
}
