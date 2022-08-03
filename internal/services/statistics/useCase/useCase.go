package statistics

import (
	"fmt"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/useCase"
	"github.com/spf13/viper"
)

type UseCase struct {
	cfg          *viper.Viper
	feedsUseCase *feeds.UseCase
}

func NewUseCaseStatistics(cfg *viper.Viper, feedsUseCase *feeds.UseCase) *UseCase {
	return &UseCase{
		cfg:          cfg,
		feedsUseCase: feedsUseCase,
	}
}

func (uc *UseCase) GatherStatistics() {
	feedsGroupByFormats := make(map[string][]int)

	getFeeds := uc.feedsUseCase.GetFeeds()

	fmt.Println(len(getFeeds))

	for _, feed := range getFeeds {
		for _, format := range feed.Formats {
			if _, ok := feedsGroupByFormats[format]; !ok {
				idFeeds := make([]int, 0)
				idFeeds = append(idFeeds, feed.Id)
				feedsGroupByFormats[format] = idFeeds
			} else {
				feedsGroupByFormats[format] = append(feedsGroupByFormats[format], feed.Id)
			}
		}
	}

	// идем за статистикой в stats-provider

}
