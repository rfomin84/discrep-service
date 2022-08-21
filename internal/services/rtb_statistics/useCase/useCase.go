package rtb_statistics

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/rfomin84/discrep-service/clients"
	useCase "github.com/rfomin84/discrep-service/internal/services/balance_history/useCase"
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

	rtbStats := make([]rtb_statistics.RtbStatistics, 0)
	feedList := u.feedsUseCase.GetFeedsWorkExternalStatistics()
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

	useCaseBalanceHistory := useCase.NewUseCaseBalanceHistory(u.cfg)
	useCaseBalanceHistory.ApprovedExternalStatistics(feedList, rtbStats, from)
}
