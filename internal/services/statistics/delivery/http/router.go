package statistics

func (d *Delivery) initRouter() {
	d.router.Use(AuthMiddleware)
	d.router.Post("/api/statistics", d.statistics)
}
