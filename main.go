package main

import (
	"golang-test1/cmd/server"
	_ "golang-test1/docs" // load API Docs files (Swagger)
	"golang-test1/pkg/config"
)

// @title Fiber Go API
// @version 1.0
// @description Fiber go web framework based REST API boilerplate
// @contact.name
// @contact.email
// @termsOfService
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:5000
// @BasePath /api
func main() {
	// setup various configuration for app
	config.LoadAllConfigs(".env")

	server.Serve()
}
