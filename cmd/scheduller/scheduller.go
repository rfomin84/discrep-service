package main

import (
	"github.com/go-co-op/gocron"
	"github.com/rfomin84/discrep-service/config"
	"github.com/rfomin84/discrep-service/internal/services/feeds/repositories/redis"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	statistics2 "github.com/rfomin84/discrep-service/internal/services/statistics/repository/long_term_storage/clickhouse"
	statisitics "github.com/rfomin84/discrep-service/internal/services/statistics/repository/temporary_storage/redis"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/useCase"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"time"
)

func main() {
	runCronJobs()
}

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)

	// get and save feeds
	s.Every(5).Minutes().SingletonMode().Do(func() {
		logger.Info("save feeds starting")
		cfg := config.GetConfig()
		repo := redis.New(cfg)
		useCaseFeeds := feeds.New(cfg, repo)
		useCaseFeeds.SaveFeeds()
		useCaseFeeds.GetFeeds()
	})

	// gather statistics from clickhouse_client
	s.Every(10).Minutes().SingletonMode().Do(func() {
		cfg := config.GetConfig()
		repo := redis.New(cfg)
		feedsUseCase := feeds.New(cfg, repo)
		tempStorageRepo := statisitics.NewTemporaryStorage(cfg)
		longTermStorageRepo := statistics2.NewLongTermStorage(cfg)

		useCaseStatistics := statistics.NewUseCaseStatistics(cfg, feedsUseCase, tempStorageRepo, longTermStorageRepo)
		useCaseStatistics.GatherStatistics()
	})

	//// gather statistics from rtb-api-provid
	//s.Every(2).Seconds().SingletonMode().Do(func() {
	//	panic("implement me")
	//})
	//
	//// calculate discrepancy
	//s.Every(2).Seconds().SingletonMode().Do(func() {
	//	panic("implement me")
	//})

	// starting cron
	s.StartBlocking()
}
