package rtb_statistics

import (
	"encoding/json"
	"time"
)

type ExternalRtbStatistics struct {
	Date        time.Time `json:"date"`
	FeedID      int       `json:"-"`
	Comment     string    `json:"comment"`
	Cost        uint64    `json:"cost"`
	Impressions uint64    `json:"impressions"`
	Clicks      uint64    `json:"clicks"`
}

type RtbStatistics struct {
	StatDate    time.Time `json:"stat_date"`
	FeedId      uint16    `json:"feed_id"`
	Country     string    `json:"country"`
	Clicks      uint64    `json:"clicks"`
	Impressions uint64    `json:"impressions"`
	Cost        uint64    `json:"cost"`
	Sign        int8      `json:"sign"`
}

func (stat *ExternalRtbStatistics) UnmarshalJSON(dataBytes []byte) error {
	type dataAlias ExternalRtbStatistics

	aliasValue := &struct {
		*dataAlias
		Date        string  `json:"date"`
		Cost        float64 `json:"cost"`
		Impressions int     `json:"impressions"`
		Clicks      int     `json:"clicks"`
	}{
		dataAlias: (*dataAlias)(stat),
	}

	if err := json.Unmarshal(dataBytes, aliasValue); err != nil {
		return err
	}

	stat.Date, _ = time.Parse("2006-01-02", aliasValue.Date)
	stat.Cost = uint64(aliasValue.Cost * 10000)
	stat.Impressions = uint64(aliasValue.Impressions)
	stat.Clicks = uint64(aliasValue.Clicks)

	return nil
}
