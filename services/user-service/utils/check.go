package utils

import (
	"fmt"
	"log"

	"github.com/abhinandpn/UnVocal/services/user-service/redis"
)

func CheckRedis() {
	err := redis.Client.Set(redis.Ctx, "name", "Abhinand", 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	value, err := redis.Client.Get(redis.Ctx, "name").Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)
}
