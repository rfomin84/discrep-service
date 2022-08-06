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
	"github.com/spf13/viper"
	"io"
	"log"
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
		fmt.Println(err.Error())
	}
}

func (uc *UseCase) FinalizeGatherStatistics() {
	startDate := carbon.Yesterday().StartOfDay().Carbon2Time()
	endDate := carbon.Now().EndOfDay().Carbon2Time()

	feedsGroupByFormats := uc.getFeeds()

	detailStatistics := uc.getStatsFromStatsProvider(feedsGroupByFormats, startDate, endDate)

	// save clickhouse
	fmt.Println(len(detailStatistics))
	uc.longTermStorageRepository.SaveStatistics(detailStatistics)
}

func (uc *UseCase) getFeeds() map[string][]int {
	feedsGroupByFormats := make(map[string][]int)

	getFeeds := uc.feedsUseCase.GetFeeds()

	fmt.Println(len(getFeeds))

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
	fmt.Println("billingType : " + billingType)
	statsProviderClient := clients.NewStatsProviederClient(uc.cfg)
	response, err := statsProviderClient.GetStatistics(startDate, endDate, billingType, timeframe, feedIds)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	statisticsStatsProvider := make([]statistics.StatisticStatsProvider, 0)
	responseBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &statisticsStatsProvider)

	if err != nil {
		log.Println(err.Error())
	}

	for _, stat := range statisticsStatsProvider {
		detailStats := statistics.DetailedFeedStatistic{
			StatDate:    stat.Date,
			FeedId:      stat.FeedId,
			BillingType: billingType,
			Country:     stat.Country,
			Clicks:      stat.Clicks,
			Impressions: stat.Impressions,
			Cost:        stat.Cost,
			Sign:        int8(1),
		}
		*stats = append(*stats, detailStats)
	}
}
