package server

// TODO  TEST!!!!

func (a *Application) RegisterRoutes() {
	r := a.Router

	r.GET("/health", a.HealthCheck)

	// Users
	userRoutes := r.Group("/users")

	userRoutes.POST("/", a.UserCreate)

	// Posts
}
