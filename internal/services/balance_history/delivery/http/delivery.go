package balance_history

import (
	"github.com/go-chi/chi/v5"
	balance_history "github.com/rfomin84/discrep-service/internal/services/balance_history/useCase"
	"github.com/spf13/viper"
)

type BalanceHistoryDelivery struct {
	cfg                   *viper.Viper
	balanceHistoryUseCase *balance_history.UseCase
	router                chi.Router
}

func NewBalanceHistoryDelivery(cfg *viper.Viper, useCase *balance_history.UseCase, router chi.Router) *BalanceHistoryDelivery {
	return &BalanceHistoryDelivery{
		cfg:                   cfg,
		balanceHistoryUseCase: useCase,
		router:                router,
	}
}

func (d *BalanceHistoryDelivery) Run() {
	d.initRouter()
}
