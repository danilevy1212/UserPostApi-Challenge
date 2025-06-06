package server

func (a *Application) RegisterRoutes() {
	r := a.Router

	r.GET("/health", a.HealthCheck)

	// Users
	userRoutes := r.Group("/users")

	userRoutes.POST("", a.UserCreate)
	userRoutes.GET("", a.UserGetAll)
	userRoutes.GET("/:id", a.UserGetByID)
	userRoutes.DELETE("/:id", a.UserDeleteByID)
	userRoutes.PUT("/:id", a.UserUpdateByID)

	// Posts
	postRoutes := r.Group("/posts")
	postRoutes.POST("", a.PostCreate)
	postRoutes.GET("", a.PostGetAll)
	postRoutes.GET("/:id", a.PostGetByID)
	postRoutes.DELETE("/:id", a.PostDeleteByID)
	postRoutes.PUT("/:id", a.PostUpdateByID)
}
