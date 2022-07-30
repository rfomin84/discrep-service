package redis

import (
	"fmt"
	"github.com/go-redis/redis/v9"
)

type Store struct {
	Client *redis.Client
}

func New(host, port, password string, DB int) *Store {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       DB,
	})

	return &Store{
		Client: client,
	}
}
