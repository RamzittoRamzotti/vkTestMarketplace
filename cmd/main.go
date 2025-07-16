package main

import (
	"github.com/joho/godotenv"
	"log"
	"vkTestMarketplace/internal/app"
	"vkTestMarketplace/internal/config"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
