package statistics

import "github.com/go-chi/chi/v5"

func (d *Delivery) initRouter() {
	d.router.Route("/api/statistics", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/", d.statistics)
	})
}
