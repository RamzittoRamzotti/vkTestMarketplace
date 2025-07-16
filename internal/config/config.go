package config

import (
	"os"
)

type Config struct {
	JWTSecret string
	DBPath    string
	Port      string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	return &Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
		DBPath:    os.Getenv("DB_PATH"),
		Port:      port,
	}
}
