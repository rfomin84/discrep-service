package balance_history

import "github.com/go-chi/chi/v5"

func (d *BalanceHistoryDelivery) initRouter() {
	d.router.Route("/api/reserve_balance", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/", d.reserveBalance)
	})
}
