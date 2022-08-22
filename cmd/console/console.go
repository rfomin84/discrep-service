package main

import (
	"fmt"
	"github.com/rfomin84/discrep-service/config"
	"github.com/rfomin84/discrep-service/internal/services/feeds/repositories/redis"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	rtb_statistics2 "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/repository/clickhouse"
	rtb_statistics "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/useCase"
	statistics2 "github.com/rfomin84/discrep-service/internal/services/statistics/repository/long_term_storage/clickhouse"
	statisitics "github.com/rfomin84/discrep-service/internal/services/statistics/repository/temporary_storage/redis"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/useCase"
	"github.com/spf13/cobra"
	"time"
)

func main() {
	fmt.Println("console")
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(gatherRtbStatistics())
	rootCmd.AddCommand(gatherStatistics())
	rootCmd.Execute()
}

func gatherRtbStatistics() *cobra.Command {
	return &cobra.Command{
		Use:   "gather-rtb-statistics [date]",
		Short: "Gather rtb statistics from external api",
		Long:  "Gather rtb statistics from external api",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			date, _ := time.Parse("2006-01-02", args[0])
			fmt.Println(date)
			cfg := config.GetConfig()
			repo := redis.New(cfg)
			feedsUseCase := feeds.New(cfg, repo)
			rtbStatisticsStorage := rtb_statistics2.NewRtbStatisticsStorage(cfg)
			useCaseRtbStatistics := rtb_statistics.NewUseCaseRtbApiStatistics(cfg, feedsUseCase, rtbStatisticsStorage)
			useCaseRtbStatistics.GatherRtbStatistics()
		},
	}
}

func gatherStatistics() *cobra.Command {
	return &cobra.Command{
		Use:   "gather-statistics [startDate] [endDate]",
		Short: "Gather rtb statistics from stats-provider",
		Long:  "Gather rtb statistics from stats-provider",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			date, _ := time.Parse("2006-01-02", args[0])
			fmt.Println(date)
			cfg := config.GetConfig()
			repo := redis.New(cfg)
			feedsUseCase := feeds.New(cfg, repo)
			tempStorageRepo := statisitics.NewTemporaryStorage(cfg)
			longTermStorageRepo := statistics2.NewLongTermStorage(cfg)
			useCaseStatistics := statistics.NewUseCaseStatistics(cfg, feedsUseCase, tempStorageRepo, longTermStorageRepo)
			useCaseStatistics.FinalizeGatherStatistics()
		},
	}
}
