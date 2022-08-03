package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/rfomin84/discrep-service/config"
	"github.com/rfomin84/discrep-service/internal/services/feeds/repositories/redis"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/useCase"
	"log"
	"time"
)

func main() {
	cfg := config.GetConfig()
	repo := redis.New(cfg)
	feedsUseCase := feeds.New(cfg, repo)
	useCaseStatistics := statistics.NewUseCaseStatistics(cfg, feedsUseCase)
	useCaseStatistics.GatherStatistics()
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
