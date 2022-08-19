package balance_history

import (
	"fmt"
	balance_history2 "github.com/rfomin84/discrep-service/internal/services/balance_history/domain"
	balance_history "github.com/rfomin84/discrep-service/internal/services/balance_history/repository/mysql"
	feeds2 "github.com/rfomin84/discrep-service/internal/services/feeds/domain"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"time"
)

type UseCase struct {
	repo *balance_history.BalanceHistoryStorage
}

func NewUseCaseBalanceHistory(cfg *viper.Viper) *UseCase {

	repo := balance_history.NewBalanceHistoryStorage(cfg)

	return &UseCase{
		repo: repo,
	}
}

func (useCase *UseCase) SaveTodayStatistics(data []statistics.DetailedFeedStatistic) {

	logger.Info(fmt.Sprintf("save today statistics. count records: %d", len(data)))
	balanceHistoryData := make([]balance_history2.BalanceHistory, 0)

	for _, stat := range data {
		balanceHistoryData = append(balanceHistoryData, balance_history2.BalanceHistory{
			FeedId:   int(stat.FeedId),
			Date:     stat.StatDate,
			Cost:     int(stat.Cost),
			Approved: false,
		})
	}
	logger.Info("save to mysql")
	useCase.repo.Save(balanceHistoryData)
}

func (useCase *UseCase) ApprovedOurStats(allFeeds []feeds2.Feed, detailStatistics []statistics.DetailedFeedStatistic, start, end time.Time) {
	balanceHistoryData := make([]balance_history2.BalanceHistory, 0)

	fmt.Println("AllFeedsCount", len(allFeeds))
	fmt.Println("CountDetailStats", len(detailStatistics))

	feedsOurStats := funk.Filter(allFeeds, func(feed feeds2.Feed) bool {
		return feed.ExternalStatistics == false
	})

	feedIds := funk.Map(feedsOurStats, func(feed feeds2.Feed) int {
		return feed.Id
	})

	fmt.Println("feedIds", len(feedIds.([]int)))
	//fmt.Println(feedIds.([]int))
	OurDetailStatistics := funk.Filter(detailStatistics, func(stat statistics.DetailedFeedStatistic) bool {
		return funk.Contains(feedIds.([]int), int(stat.FeedId))
	})

	fmt.Println("OurDetailStats", len(OurDetailStatistics.([]statistics.DetailedFeedStatistic)))

	for _, stat := range OurDetailStatistics.([]statistics.DetailedFeedStatistic) {
		balanceHistoryData = append(balanceHistoryData, balance_history2.BalanceHistory{
			FeedId:   int(stat.FeedId),
			Date:     stat.StatDate,
			Cost:     int(stat.Cost),
			Approved: true,
		})
	}

	useCase.repo.DeleteStatisticsByFeedIds(feedIds.([]int), start, end, false)
	useCase.repo.Save(balanceHistoryData)
}

func (useCase *UseCase) DeleteNotApprovedStatistics(start, end time.Time) {
	useCase.repo.DeleteNotApproveStatistics(start, end)
}

func (useCase *UseCase) ReservedFeedBalance() ([]balance_history2.ReservedBalance, error) {
	return useCase.repo.GetReserveFeedBalance()
}
