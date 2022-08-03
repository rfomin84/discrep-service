package statistics

type DetailedFeedStatistic struct {
	StatDate    string  `json:"date"`
	FeedId      int     `json:"feed_id"`
	BillingType string  `json:"-"`
	Country     string  `json:"country"`
	Clicks      int     `json:"clicks"`
	Impressions int     `json:"impressions"`
	Cost        float64 `json:"cost"`
	Sign        int8    `json:"-"`
}

type

