package rtb_statistics

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/rfomin84/discrep-service/clients"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	rtb_statistics "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/domain"
	rtb_statistics2 "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/repository"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"github.com/spf13/viper"
	"sync"
)

type UseCase struct {
	cfg                  *viper.Viper
	RtbApiProviderClient *clients.RtbApiProvidClient
	feedsUseCase         *feeds.UseCase
	storage              rtb_statistics2.RtbStatisticStorageInterface
}

func NewUseCaseRtbApiStatistics(cfg *viper.Viper, feedUseCase *feeds.UseCase, storage rtb_statistics2.RtbStatisticStorageInterface) *UseCase {
	rtbApiProviderClient := clients.NewRtbApiProvidClient(*cfg)
	return &UseCase{
		cfg:                  cfg,
		RtbApiProviderClient: rtbApiProviderClient,
		feedsUseCase:         feedUseCase,
		storage:              storage,
	}
}

func (u *UseCase) GatherRtbStatistics() {
	from := carbon.Yesterday().Carbon2Time()
	//to := carbon.Yesterday().Carbon2Time()
	rtbStats := make([]rtb_statistics.RtbStatistics, 0)
	feedList := u.feedsUseCase.GetFeeds()
	fmt.Println(feedList)
	countWorkers := 20
	taskCh := make(chan Task, countWorkers)
	var wg sync.WaitGroup

	for i := 0; i < countWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskCh {
				worker(task, u.RtbApiProviderClient, &rtbStats)
			}
		}()
	}

	go func() {
		for _, feed := range feedList {
			logger.Info(fmt.Sprintf("task create for feed %d with rtb_provider_id: %d", feed.Id, feed.RtbApiProviderId))
			taskCh <- Task{
				FeedID:           feed.Id,
				Date:             from,
				RtbApiProviderId: feed.RtbApiProviderId,
			}
		}
		close(taskCh)

	}()

	wg.Wait()

	logger.Info(fmt.Sprintf("count rtbStats : %d", len(rtbStats)))

	u.storage.SaveRtbStatistics(rtbStats)

	//for _, feed := range feedList {
	//	if feed.Id > 500 {
	//		break
	//	}
	//	logger.Info(fmt.Sprintf("rtb_api_provider_id : %d", feed.RtbApiProviderId))
	//	response, err := u.RtbApiProviderClient.GetStatistics(from, to, strconv.Itoa(feed.RtbApiProviderId))
	//	if err != nil {
	//		logger.Warning("Error get external statistics: " + err.Error())
	//		continue
	//	}
	//	if response.StatusCode != 200 {
	//		logger.Warning(fmt.Sprintf("Not get rtb statistics for feed %d", feed.Id))
	//		continue
	//	}
	//
	//	// save statistics
	//	responseBody, _ := io.ReadAll(response.Body)
	//	var extRtbStat rtb_statistics.ExternalRtbStatistics
	//
	//	err = json.Unmarshal(responseBody, &extRtbStat)
	//	if err != nil {
	//		logger.Error(err.Error())
	//	}
	//	rtbStats = append(rtbStats, rtb_statistics.RtbStatistics{
	//		StatDate:    extRtbStat.Date,
	//		FeedId:      uint16(feed.Id),
	//		Country:     "",
	//		Clicks:      extRtbStat.Clicks,
	//		Impressions: extRtbStat.Impressions,
	//		Cost:        extRtbStat.Cost,
	//		Sign:        int8(1),
	//	})
	//}
	//
	//u.storage.SaveRtbStatistics(rtbStats)
}
