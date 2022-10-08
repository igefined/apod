package main

import (
	"github.com/igilgyrg/betera-test/internal/app"
)

// @title betera
// @version 1.0
// @description API Server betera test service

// @host http://localhost:3000
// @BasePath /
func main() {
	application := app.New()
	application.Start()
}
