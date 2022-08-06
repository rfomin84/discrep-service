package statistics

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
	"github.com/rfomin84/discrep-service/pkg/store/clickhouse_client"
	"github.com/spf13/viper"
	"log"
)

type LongTermStorage struct {
	conn driver.Conn
}

func NewLongTermStorage(cfg *viper.Viper) *LongTermStorage {
	conn, err := clickhouse_client.NewClickhouseClient(
		context.Background(),
		cfg.GetString("CLICKHOUSE_HOST"),
		cfg.GetString("CLICKHOUSE_PORT"),
		cfg.GetString("CLICKHOUSE_USERNAME"),
		cfg.GetString("CLICKHOUSE_PASSWORD"),
		cfg.GetString("CLICKHOUSE_DATABASE"),
	)

	if err != nil {
		log.Println("error connect", err.Error())
	}

	return &LongTermStorage{
		conn: conn,
	}
}

func (repo *LongTermStorage) SaveStatistics(stats []statistics.DetailedFeedStatistic) {

	batch, err := repo.conn.PrepareBatch(
		context.Background(),
		"INSERT INTO detailed_feed_statistics (StatDate, FeedId, BillingType, Country, Clicks, Impressions, Cost, Sign)",
	)
	if err != nil {
		log.Println(err.Error())
	}

	for _, item := range stats {
		fmt.Println(item)
		if err := batch.Append(
			item.StatDate,
			uint16(item.FeedId),
			item.BillingType,
			item.Country,
			uint64(item.Clicks),
			uint64(item.Impressions),
			uint64(item.Cost),
			item.Sign,
		); err != nil {
			log.Println(err.Error())
		}
	}

	if err := batch.Send(); err != nil {
		log.Println(err.Error())
	}
}
