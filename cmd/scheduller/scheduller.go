package main

import (
	"context"
	"fmt"
	"github.com/rfomin84/discrep-service/config"
	"github.com/rfomin84/discrep-service/pkg/store/redis"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println(cfg.Get("API_TOKEN"))
	fmt.Println(cfg.Get("REDIS_HOST"))

	redisStore := redis.New(
		cfg.GetString("REDIS_HOST"),
		cfg.GetString("REDIS_PORT"),
		cfg.GetString("REDIS_PASSWORD"),
		cfg.GetInt("REDIS_DB"),
	)
	ctx := context.Background()
	err := redisStore.Client.Set(ctx, "id1234", "test", 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	val, err := redisStore.Client.Get(ctx, "id1234").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}
