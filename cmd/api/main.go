package main

import (
	"fmt"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/server"
)

func main() {
	app := server.New()

	app.RegisterMiddleware()
	app.RegisterRoutes()

	if err := app.Serve(app.Config.Port); err != nil {
		fmt.Println("Fatal error occured", err)
	}
}
