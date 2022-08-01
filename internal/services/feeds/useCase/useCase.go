package feeds

import (
	"context"
	"github.com/rfomin84/discrep-service/clients"
	feeds2 "github.com/rfomin84/discrep-service/internal/services/feeds/domain"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/repositories"
	"github.com/spf13/viper"
	"io"
	"log"
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
		log.Println("error get feeds from tc3")
		return
	}
	defer feedsData.Body.Close()
	bytes, err := io.ReadAll(feedsData.Body)
	if err != nil {
		log.Println("error", err.Error())
	}

	// сохранить в хранилище данные
	if err = uc.repository.Save(context.Background(), "feeds", bytes); err != nil {
		log.Println("error save feeds to storage", err.Error())
	}
}

func (uc *UseCase) GetFeeds() []feeds2.Feed {
	// получить дынные из хранилища
	feedsAll, err := uc.repository.Get(context.Background(), "feeds")
	if err != nil {
		return nil
	}
	return feedsAll
}
