package config

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// LoadAllConfigs set various configs
func LoadAllConfigs(envFile string) {

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("can't load .env file. error: %v", err)
	}

	LoadApp()
	LoadDBCfg()
}

// FiberConfig func for configuration Fiber app.
func FiberConfig() fiber.Config {

	// Return Fiber configuration.
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(AppCfg().ReadTimeout),
	}
}
