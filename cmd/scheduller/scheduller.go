package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/rfomin84/discrep-service/config"
	"github.com/rfomin84/discrep-service/internal/services/feeds/repositories/redis"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	rtb_statistics2 "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/repository/clickhouse"
	rtb_statistics "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/useCase"
	"log"
	"time"
)

func main() {
	cfg := config.GetConfig()
	repo := redis.New(cfg)
	feedsUseCase := feeds.New(cfg, repo)

	// gather statistics
	//tempStorageRepo := statisitics.NewTemporaryStorage(cfg)
	//longTermStorageRepo := statistics2.NewLongTermStorage(cfg)
	//fmt.Println(longTermStorageRepo)
	//useCaseStatistics := statistics.NewUseCaseStatistics(cfg, feedsUseCase, tempStorageRepo, longTermStorageRepo)
	//useCaseStatistics.GatherStatistics()
	//useCaseStatistics.FinalizeGatherStatistics()

	// gather rtb statistics
	rtbStatisticsStorage := rtb_statistics2.NewRtbStatisticsStorage(cfg)
	useCaseRtbStatistics := rtb_statistics.NewUseCaseRtbApiStatistics(cfg, feedsUseCase, rtbStatisticsStorage)
	useCaseRtbStatistics.GatherRtbStatistics()

	//runCronJobs()
}

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)

	// get and save feeds
	s.Every(5).Minutes().SingletonMode().Do(func() {
		log.Println("save feeds starting")
		cfg := config.GetConfig()
		repo := redis.New(cfg)
		useCaseFeeds := feeds.New(cfg, repo)
		useCaseFeeds.SaveFeeds()
		fmt.Println(useCaseFeeds.GetFeeds())
	})

	// gather statistics from clickhouse_client
	//s.Every(10).Seconds().SingletonMode().Do(func() {
	//	cfg := config.GetConfig()
	//	repo := redis.New(cfg)
	//	useCaseFeeds := feeds.New(cfg, repo)
	//	fmt.Println(useCaseFeeds.GetFeeds())
	//})

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
