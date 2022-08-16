package rtb_statistics

import "time"

type Task struct {
	Date             time.Time
	FeedID           int
	RtbApiProviderId int
}
