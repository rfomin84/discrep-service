package balance_history

import (
	"fmt"
	balance_history2 "github.com/rfomin84/discrep-service/internal/services/balance_history/domain"
	balance_history "github.com/rfomin84/discrep-service/internal/services/balance_history/repository/mysql"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"github.com/spf13/viper"
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

func (useCase *UseCase) DeleteNotApprovedStatistics(start, end time.Time) {
	useCase.repo.DeleteNotApproveStatistics(start, end)
}

func (useCase *UseCase) ReservedFeedBalance() ([]balance_history2.ReservedBalance, error) {
	return useCase.repo.GetReserveFeedBalance()
}
