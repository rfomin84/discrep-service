package rtb_statistics

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	rtb_statistics "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/domain"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"github.com/rfomin84/discrep-service/pkg/store/clickhouse_client"
	"github.com/spf13/viper"
)

type RtbStatisticsStorage struct {
	conn driver.Conn
}

func NewRtbStatisticsStorage(cfg *viper.Viper) *RtbStatisticsStorage {
	conn, err := clickhouse_client.NewClickhouseClient(
		context.Background(),
		cfg.GetString("CLICKHOUSE_HOST"),
		cfg.GetString("CLICKHOUSE_PORT"),
		cfg.GetString("CLICKHOUSE_USERNAME"),
		cfg.GetString("CLICKHOUSE_PASSWORD"),
		cfg.GetString("CLICKHOUSE_DATABASE"),
	)

	if err != nil {
		logger.Error("error connect" + err.Error())
	}

	return &RtbStatisticsStorage{
		conn: conn,
	}
}

func (repo *RtbStatisticsStorage) SaveRtbStatistics(stats []rtb_statistics.RtbStatistics) {
	batch, err := repo.conn.PrepareBatch(
		context.Background(),
		"INSERT INTO rtb_statistics (StatDate, FeedId, Country, Clicks, Impressions, Cost, Sign)",
	)
	if err != nil {
		logger.Error(err.Error())
	}

	for _, item := range stats {
		if err := batch.Append(
			item.StatDate,
			item.FeedId,
			item.Country,
			item.Clicks,
			item.Impressions,
			item.Cost,
			item.Sign,
		); err != nil {
			logger.Error(err.Error())
		}
	}

	if err := batch.Send(); err != nil {
		logger.Error(err.Error())
	}
}
