package statistics

import (
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
	"time"
)

type LongTermStorageInterface interface {
	SaveStatistics(stats []statistics.DetailedFeedStatistic)
	GetStatistics(startDate, endDate time.Time, feedIds []uint16) []statistics.DetailedFeedStatistic
}
