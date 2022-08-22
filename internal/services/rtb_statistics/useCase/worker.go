package rtb_statistics

import (
	"encoding/json"
	"fmt"
	"github.com/rfomin84/discrep-service/clients"
	rtb_statistics "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/domain"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"io"
	"strconv"
)

func worker(task Task, client *clients.RtbApiProvidClient, stats *[]rtb_statistics.RtbStatistics) {
	response, err := client.GetStatistics(task.Date, task.Date, strconv.Itoa(task.RtbApiProviderId))
	if err != nil {
		logger.Warning("Error get external statistics: " + err.Error())
		return
	}
	if response.StatusCode != 200 {
		logger.Warning(fmt.Sprintf("Not get rtb statistics for feed %d", task.FeedID))
		return
	}

	responseBody, _ := io.ReadAll(response.Body)
	var extRtbStat rtb_statistics.ExternalRtbStatistics

	err = json.Unmarshal(responseBody, &extRtbStat)
	if err != nil {
		logger.Error(err.Error())
	}
	extRtbStat.FeedID = task.FeedID
	logger.Info(fmt.Sprintf("feedId : %d, result: %v", task.FeedID, extRtbStat))
	*stats = append(*stats, rtb_statistics.RtbStatistics{
		StatDate:    extRtbStat.Date,
		FeedId:      uint16(extRtbStat.FeedID),
		Country:     "",
		Clicks:      extRtbStat.Clicks,
		Impressions: extRtbStat.Impressions,
		Cost:        extRtbStat.Cost,
		Sign:        int8(1),
	})
}
