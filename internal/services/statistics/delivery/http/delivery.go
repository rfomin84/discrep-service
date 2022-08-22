package statistics

import (
	"github.com/go-chi/chi/v5"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/useCase"
	"github.com/spf13/viper"
)

type Delivery struct {
	cfg              *viper.Viper
	statisticUseCase *statistics.UseCase
	router           chi.Router
}

func NewHttpStatisticDelivery(cfg *viper.Viper, useCase *statistics.UseCase, router chi.Router) *Delivery {
	return &Delivery{
		cfg:              cfg,
		statisticUseCase: useCase,
		router:           router,
	}
}

func (d *Delivery) Run() {
	d.initRouter()
}
