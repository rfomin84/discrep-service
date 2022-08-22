package statistics

import (
	"encoding/json"
	"time"
)

type DetailedFeedStatistic struct {
	StatDate    time.Time `json:"date"`
	FeedId      uint16    `json:"feed_id"`
	BillingType string    `json:"-"`
	Country     string    `json:"country"`
	Clicks      uint64    `json:"clicks"`
	Impressions uint64    `json:"impressions"`
	Cost        uint64    `json:"cost"`
	Sign        int8      `json:"-"`
}

type StatisticStatsProvider struct {
	FeedId      int       `json:"feed_id"`
	Date        time.Time `json:"date"`
	Country     string    `json:"country"`
	Impressions int       `json:"impressions"`
	Clicks      int       `json:"clicks"`
	Cost        float64   `json:"cost"`
}

func (d StatisticStatsProvider) MarshalJSON() ([]byte, error) {
	type dataAlias StatisticStatsProvider

	aliasValue := struct {
		dataAlias
		Date string `json:"date"`
	}{
		dataAlias: dataAlias(d),
		Date:      d.Date.Format("2006-01-02 15:04:05"),
	}

	return json.Marshal(aliasValue)
}

func (d *StatisticStatsProvider) UnmarshalJSON(dataBytes []byte) error {
	type dataAlias StatisticStatsProvider

	aliasValue := &struct {
		*dataAlias
		Date string `json:"date"`
	}{
		dataAlias: (*dataAlias)(d),
	}

	if err := json.Unmarshal(dataBytes, aliasValue); err != nil {
		return err
	}

	d.Date, _ = time.Parse("2006-01-02 15:04:05", aliasValue.Date)

	return nil
}
