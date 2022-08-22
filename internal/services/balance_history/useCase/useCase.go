package balance_history

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	balance_history2 "github.com/rfomin84/discrep-service/internal/services/balance_history/domain"
	balance_history "github.com/rfomin84/discrep-service/internal/services/balance_history/repository/mysql"
	feeds2 "github.com/rfomin84/discrep-service/internal/services/feeds/domain"
	rtb_statistics "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/domain"
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

func (useCase *UseCase) ApprovedOurStats(feedsWorkOurStat []feeds2.Feed, detailStatistics []statistics.DetailedFeedStatistic, start, end time.Time) {
	balanceHistoryData := make([]balance_history2.BalanceHistory, 0)

	feedIds := funk.Map(feedsWorkOurStat, func(feed feeds2.Feed) int {
		return feed.Id
	})

	OurDetailStatistics := funk.Filter(detailStatistics, func(stat statistics.DetailedFeedStatistic) bool {
		return funk.Contains(feedIds.([]int), int(stat.FeedId))
	})

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

func (useCase *UseCase) ApprovedExternalStatistics(feedsWorkExternalStats []feeds2.Feed, rtbStats []rtb_statistics.RtbStatistics, date time.Time) {

	feedIds := funk.Map(feedsWorkExternalStats, func(feed feeds2.Feed) int {
		return feed.Id
	})
	carbonEl := carbon.Time2Carbon(date)
	start := carbonEl.StartOfDay().Carbon2Time()
	end := carbonEl.EndOfDay().Carbon2Time()
	useCase.repo.DeleteStatisticsByFeedIds(feedIds.([]int), start, end, false)

	balanceHistoryData := make([]balance_history2.BalanceHistory, 0)

	for _, stat := range rtbStats {
		balanceHistoryData = append(balanceHistoryData, balance_history2.BalanceHistory{
			FeedId:   int(stat.FeedId),
			Date:     stat.StatDate,
			Cost:     int(stat.Cost),
			Approved: true,
		})
	}

	useCase.repo.Save(balanceHistoryData)
}

func (useCase *UseCase) DeleteNotApprovedStatistics(start, end time.Time) {
	useCase.repo.DeleteNotApproveStatistics(start, end)
}

func (useCase *UseCase) ReservedFeedBalance() ([]balance_history2.ReservedBalance, error) {
	return useCase.repo.GetReserveFeedBalance()
}
