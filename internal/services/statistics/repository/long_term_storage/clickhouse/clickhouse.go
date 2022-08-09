package statistics

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
	"github.com/rfomin84/discrep-service/pkg/store/clickhouse_client"
	"github.com/spf13/viper"
	"log"
	"time"
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

func (repo *LongTermStorage) GetStatistics(startDate, endDate time.Time, feedIds []uint16) []statistics.DetailedFeedStatistic {

	stats := make([]statistics.DetailedFeedStatistic, 0)

	rows := "SELECT StatDate, FeedId, BillingType, Country, SUM(Clicks) as Clicks, SUM(Impressions) as Impressions, SUM(Cost) as Cost " +
		"FROM detailed_feed_statistics " +
		"WHERE StatDate >= toDateTime(@startDate) AND StatDate <= toDateTime(@endDate) AND FeedId IN (@feedIds) " +
		"GROUP BY StatDate, FeedId, BillingType, Country ORDER BY StatDate"

	if err := repo.conn.Select(context.Background(), &stats, rows,
		clickhouse.Named("startDate", startDate),
		clickhouse.Named("endDate", endDate),
		clickhouse.Named("feedIds", feedIds),
	); err != nil {
		log.Println(err.Error())
		return stats
	}

	return stats
}
