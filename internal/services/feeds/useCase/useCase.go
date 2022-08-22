package feeds

import (
	"context"
	"github.com/rfomin84/discrep-service/clients"
	feeds2 "github.com/rfomin84/discrep-service/internal/services/feeds/domain"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/repositories"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"io"
)

type UseCase struct {
	cfg        *viper.Viper
	repository feeds.StoreInterface
}

func New(cfg *viper.Viper, repo feeds.StoreInterface) *UseCase {
	return &UseCase{
		cfg:        cfg,
		repository: repo,
	}
}

func (uc *UseCase) SaveFeeds() {
	// получить дынные из tc3
	tc3Client := clients.New(uc.cfg)
	feedsData, err := tc3Client.GetFeeds()
	if err != nil {
		logger.Error("error get feeds from tc3 : " + err.Error())
		return
	}
	defer feedsData.Body.Close()
	bytes, err := io.ReadAll(feedsData.Body)
	if err != nil {
		logger.Error("Error : " + err.Error())
	}

	// сохранить в хранилище данные
	if err = uc.repository.Save(context.Background(), "feeds", bytes); err != nil {
		logger.Error("error save feeds to storage" + err.Error())
	}
}

func (uc *UseCase) GetFeeds() []feeds2.Feed {
	// получить дынные из хранилища
	feedsAll, err := uc.repository.Get(context.Background(), "feeds")
	if err != nil {
		logger.Warning(err.Error())
		return nil
	}
	return feedsAll
}

func (uc *UseCase) GetFeedsWorkOurStatistics() []feeds2.Feed {
	allFeeds := uc.GetFeeds()
	feedsOurStats := funk.Filter(allFeeds, func(feed feeds2.Feed) bool {
		return feed.ExternalStatistics == false
	})

	return feedsOurStats.([]feeds2.Feed)
}

func (uc *UseCase) GetFeedsWorkExternalStatistics() []feeds2.Feed {
	allFeeds := uc.GetFeeds()
	feedsOurStats := funk.Filter(allFeeds, func(feed feeds2.Feed) bool {
		return feed.ExternalStatistics == true
	})

	return feedsOurStats.([]feeds2.Feed)
}
