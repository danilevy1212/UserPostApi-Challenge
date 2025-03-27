package server

// TODO  TEST!!!!

func (a *Application) RegisterRoutes() {
	r := a.router

	r.GET("/health", a.HealthCheck)
}
