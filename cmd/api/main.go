package main

import (
	"fmt"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/server"
)

func main() {
	app := server.New()

	app.RegisterMiddleware()
	app.RegisterRoutes()

	// TODO  Hardcoded port for now
	if err := app.Serve(3000); err != nil {
		fmt.Println("Fatal error occured", err)
	}
}
