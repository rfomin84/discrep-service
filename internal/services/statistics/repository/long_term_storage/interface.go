package statistics

import statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"

type LongTermStorageInterface interface {
	SaveStatistics(stats []statistics.DetailedFeedStatistic)
}
