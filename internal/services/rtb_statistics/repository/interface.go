package rtb_statistics

import rtb_statistics "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/domain"

type RtbStatisticStorageInterface interface {
	SaveRtbStatistics(stats []rtb_statistics.RtbStatistics)
}
