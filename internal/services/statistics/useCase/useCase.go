package statistics

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/rfomin84/discrep-service/clients"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
	statistics3 "github.com/rfomin84/discrep-service/internal/services/statistics/repository/long_term_storage"
	statistics2 "github.com/rfomin84/discrep-service/internal/services/statistics/repository/temporary_storage"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"github.com/spf13/viper"
	"io"
	"sync"
	"time"
)

type UseCase struct {
	cfg                        *viper.Viper
	feedsUseCase               *feeds.UseCase
	temporaryStorageRepository statistics2.TemporaryStorageInterface
	longTermStorageRepository  statistics3.LongTermStorageInterface
}

func NewUseCaseStatistics(cfg *viper.Viper, feedsUseCase *feeds.UseCase, temporaryRepo statistics2.TemporaryStorageInterface, longTermRepo statistics3.LongTermStorageInterface) *UseCase {
	return &UseCase{
		cfg:                        cfg,
		feedsUseCase:               feedsUseCase,
		temporaryStorageRepository: temporaryRepo,
		longTermStorageRepository:  longTermRepo,
	}
}

func (uc *UseCase) GatherStatistics() {

	startDate := carbon.Yesterday().StartOfDay().Carbon2Time()
	endDate := carbon.Now().EndOfDay().Carbon2Time()

	feedsGroupByFormats := uc.getFeeds()

	detailStatistics := uc.getStatsFromStatsProvider(feedsGroupByFormats, startDate, endDate)

	err := uc.temporaryStorageRepository.SaveStatistics(context.Background(), detailStatistics)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (uc *UseCase) FinalizeGatherStatistics() {
	startDate := carbon.Yesterday().StartOfDay().Carbon2Time()
	endDate := carbon.Now().EndOfDay().Carbon2Time()

	feedsGroupByFormats := uc.getFeeds()

	detailStatistics := uc.getStatsFromStatsProvider(feedsGroupByFormats, startDate, endDate)

	// save clickhouse
	uc.longTermStorageRepository.SaveStatistics(detailStatistics)
}

func (uc *UseCase) GetStatistics(startDate, endDate string, feedIds []int) []statistics.DetailedFeedStatistic {
	start, _ := time.Parse("2006-01-02 15:04:05", startDate)
	end, _ := time.Parse("2006-01-02 15:04:05", endDate)

	feedListId := make([]uint16, 0)

	for _, feedId := range feedIds {
		feedListId = append(feedListId, uint16(feedId))
	}

	return uc.longTermStorageRepository.GetStatistics(start, end, feedListId)
}

func (uc *UseCase) getFeeds() map[string][]int {
	feedsGroupByFormats := make(map[string][]int)

	getFeeds := uc.feedsUseCase.GetFeeds()
	logger.Debug(fmt.Sprintf("Count get feeds: %d", len(getFeeds)))

	for _, feed := range getFeeds {
		/** TODO: не учитывается формирование формата dsp + billing_type */
		for _, format := range feed.Formats {
			if _, ok := feedsGroupByFormats[format]; !ok {
				idFeeds := make([]int, 0)
				idFeeds = append(idFeeds, feed.Id)
				feedsGroupByFormats[format] = idFeeds
			} else {
				feedsGroupByFormats[format] = append(feedsGroupByFormats[format], feed.Id)
			}
		}
	}

	return feedsGroupByFormats
}

func (uc *UseCase) getStatsFromStatsProvider(feedsGroupByFormats map[string][]int, startDate, endDate time.Time) []statistics.DetailedFeedStatistic {
	detailStatistics := make([]statistics.DetailedFeedStatistic, 0)
	var wg sync.WaitGroup

	for billingType, feedIds := range feedsGroupByFormats {
		wg.Add(1)
		go uc.getStatisticByBillingType(&wg, &detailStatistics, startDate, endDate, billingType, "hour", feedIds)
	}

	wg.Wait()

	return detailStatistics
}

func (uc *UseCase) getStatisticByBillingType(wg *sync.WaitGroup, stats *[]statistics.DetailedFeedStatistic, startDate, endDate time.Time, billingType, timeframe string, feedIds []int) {
	defer wg.Done()
	logger.Info(fmt.Sprintf("billingType : %s", billingType))
	statsProviderClient := clients.NewStatsProviederClient(uc.cfg)
	response, err := statsProviderClient.GetStatistics(startDate, endDate, billingType, timeframe, feedIds)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	statisticsStatsProvider := make([]statistics.StatisticStatsProvider, 0)
	responseBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &statisticsStatsProvider)

	if err != nil {
		logger.Error(err.Error())
	}

	for _, stat := range statisticsStatsProvider {
		detailStats := statistics.DetailedFeedStatistic{
			StatDate:    stat.Date,
			FeedId:      uint16(stat.FeedId),
			BillingType: billingType,
			Country:     stat.Country,
			Clicks:      uint64(stat.Clicks),
			Impressions: uint64(stat.Impressions),
			Cost:        uint64(stat.Cost),
			Sign:        int8(1),
		}
		*stats = append(*stats, detailStats)
	}
}
