package statisitics

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
	"github.com/rfomin84/discrep-service/pkg/store/redis_client"
	"github.com/spf13/viper"
)

type TemporaryStorage struct {
	client *redis.Client
}

func NewTemporaryStorage(cfg *viper.Viper) *TemporaryStorage {
	redisStore := redis_client.New(
		cfg.GetString("REDIS_HOST"),
		cfg.GetString("REDIS_PORT"),
		cfg.GetString("REDIS_PASSWORD"),
		cfg.GetInt("REDIS_DB"),
	)
	return &TemporaryStorage{
		client: redisStore,
	}
}

func (ts *TemporaryStorage) SaveStatistics(ctx context.Context, stats []statistics.DetailedFeedStatistic) error {

	byteData, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return ts.client.Set(ctx, "statistics", byteData, 0).Err()
}
