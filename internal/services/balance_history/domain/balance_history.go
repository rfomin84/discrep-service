package balance_history

import "time"

type BalanceHistory struct {
	ID       int
	FeedId   int
	Date     time.Time
	Cost     int
	Approved bool
}

type ReservedBalance struct {
	FeedId int `json:"feed_id"`
	Cost   int `json:"cost"`
}
