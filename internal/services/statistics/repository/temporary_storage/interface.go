package statistics

import (
	"context"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
)

type TemporaryStorageInterface interface {
	SaveStatistics(ctx context.Context, stats []statistics.DetailedFeedStatistic) error
}
