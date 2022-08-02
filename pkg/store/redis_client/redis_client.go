package redis_client

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
)

func New(host, port, password string, DB int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       DB,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		log.Fatal("no connect with redis_client")
	}

	return client
}
