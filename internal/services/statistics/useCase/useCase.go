package statistics

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/rfomin84/discrep-service/clients"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
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
}

func NewUseCaseStatistics(cfg *viper.Viper, feedsUseCase *feeds.UseCase, repo statistics2.TemporaryStorageInterface) *UseCase {
	return &UseCase{
		cfg:                        cfg,
		feedsUseCase:               feedsUseCase,
		temporaryStorageRepository: repo,
	}
}

func (uc *UseCase) GatherStatistics() {
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

	// идем за статистикой в stats-provider
	detailStatistics := make([]statistics.DetailedFeedStatistic, 0)
	var wg sync.WaitGroup

	startDate := carbon.Yesterday().StartOfDay().Carbon2Time()
	endDate := carbon.Now().EndOfDay().Carbon2Time()

	for billingType, feedIds := range feedsGroupByFormats {
		wg.Add(1)
		go uc.getStatisticByBillingType(&wg, &detailStatistics, startDate, endDate, billingType, "hour", feedIds)
	}

	wg.Wait()

	fmt.Println(len(detailStatistics))

	err := uc.temporaryStorageRepository.SaveStatistics(context.Background(), detailStatistics)
	if err != nil {
		fmt.Println(err.Error())
	}

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
		}
		*stats = append(*stats, detailStats)
	}
}
