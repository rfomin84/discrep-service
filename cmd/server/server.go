package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rfomin84/discrep-service/config"
	"github.com/rfomin84/discrep-service/internal/services/feeds/repositories/redis"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/delivery/http"
	statistics3 "github.com/rfomin84/discrep-service/internal/services/statistics/repository/long_term_storage/clickhouse"
	statisitics "github.com/rfomin84/discrep-service/internal/services/statistics/repository/temporary_storage/redis"
	statistics2 "github.com/rfomin84/discrep-service/internal/services/statistics/useCase"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()
	router := chi.NewRouter()
	
	repo := redis.New(cfg)

	feedsUseCase := feeds.New(cfg, repo)
	tempStorageRepo := statisitics.NewTemporaryStorage(cfg)
	longTermStorageRepo := statistics3.NewLongTermStorage(cfg)
	useCaseStatistics := statistics2.NewUseCaseStatistics(cfg, feedsUseCase, tempStorageRepo, longTermStorageRepo)
	statisticsDeliveryStatistics := statistics.NewHttpStatisticDelivery(cfg, useCaseStatistics, router)
	statisticsDeliveryStatistics.Run()

	if err := http.ListenAndServe(":"+cfg.GetString("APP_PORT"), router); err != nil {
		log.Fatalln("server error starting", err.Error())
	}
}
