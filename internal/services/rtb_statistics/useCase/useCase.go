package rtb_statistics

import (
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/rfomin84/discrep-service/clients"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	rtb_statistics "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/domain"
	rtb_statistics2 "github.com/rfomin84/discrep-service/internal/services/rtb_statistics/repository"
	"github.com/spf13/viper"
	"io"
	"log"
	"strconv"
)

type UseCase struct {
	cfg                  *viper.Viper
	RtbApiProviderClient *clients.RtbApiProvidClient
	feedsUseCase         *feeds.UseCase
	storage              rtb_statistics2.RtbStatisticStorageInterface
}

func NewUseCaseRtbApiStatistics(cfg *viper.Viper, feedUseCase *feeds.UseCase, storage rtb_statistics2.RtbStatisticStorageInterface) *UseCase {
	rtbApiProviderClient := clients.NewRtbApiProvidClient(*cfg)
	return &UseCase{
		cfg:                  cfg,
		RtbApiProviderClient: rtbApiProviderClient,
		feedsUseCase:         feedUseCase,
		storage:              storage,
	}
}

func (u *UseCase) GatherRtbStatistics() {
	from := carbon.Yesterday().Carbon2Time()
	to := carbon.Yesterday().Carbon2Time()
	rtbStats := make([]rtb_statistics.RtbStatistics, 0)
	feedList := u.feedsUseCase.GetFeeds()

	for _, feed := range feedList {
		if feed.Id > 500 {
			break
		}
		fmt.Println("rtb_api_provider_id", feed.RtbApiProviderId)
		response, err := u.RtbApiProviderClient.GetStatistics(from, to, strconv.Itoa(feed.RtbApiProviderId))
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if response.StatusCode != 200 {
			log.Println(fmt.Sprintf("Not get rtb statistics for feed %d", feed.Id))
			continue
		}

		// save statistics
		responseBody, _ := io.ReadAll(response.Body)
		var extRtbStat rtb_statistics.ExternalRtbStatistics
		fmt.Println(string(responseBody))
		err = json.Unmarshal(responseBody, &extRtbStat)
		if err != nil {
			log.Println(err.Error())
		}
		rtbStats = append(rtbStats, rtb_statistics.RtbStatistics{
			StatDate:    extRtbStat.Date,
			FeedId:      uint16(feed.Id),
			Country:     "",
			Clicks:      extRtbStat.Clicks,
			Impressions: extRtbStat.Impressions,
			Cost:        extRtbStat.Cost,
			Sign:        int8(1),
		})
	}
	fmt.Println(rtbStats)

	u.storage.SaveRtbStatistics(rtbStats)
}
