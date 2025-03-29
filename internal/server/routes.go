package server

// TODO  TEST!!!!

func (a *Application) RegisterRoutes() {
	r := a.Router

	r.GET("/health", a.HealthCheck)
}
