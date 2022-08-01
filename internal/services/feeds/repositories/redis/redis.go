package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	feeds2 "github.com/rfomin84/discrep-service/internal/services/feeds/domain"
	"github.com/rfomin84/discrep-service/pkg/store/redis_client"
	"github.com/spf13/viper"
)

type Storage struct {
	client *redis.Client
}

func New(cfg *viper.Viper) *Storage {
	redisStore := redis_client.New(
		cfg.GetString("REDIS_HOST"),
		cfg.GetString("REDIS_PORT"),
		cfg.GetString("REDIS_PASSWORD"),
		cfg.GetInt("REDIS_DB"),
	)
	return &Storage{
		client: redisStore.Client,
	}
}

func (storage *Storage) Save(ctx context.Context, key string, data interface{}) error {
	return storage.client.Set(ctx, key, data, 0).Err()
}

func (storage *Storage) Get(ctx context.Context, key string) ([]feeds2.Feed, error) {
	val, err := storage.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	feeds := make([]feeds2.Feed, 0)

	err = json.Unmarshal([]byte(val), &feeds)

	if err != nil {
		return nil, err
	}

	return feeds, nil
}
