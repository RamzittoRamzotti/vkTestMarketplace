package app

import (
	"log"
	"net/http"

	"vkTestMarketplace/internal/config"
	"vkTestMarketplace/internal/http-server/handlers"
	"vkTestMarketplace/internal/http-server/router"
	"vkTestMarketplace/internal/storage/sqlite"
)

type App struct {
	Cfg    *config.Config
	Router http.Handler
}

func New(cfg *config.Config) (*App, error) {
	store, err := sqlite.New(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	authHandlers := &handlers.AuthHandlers{
		Users:     store,
		JWTSecret: cfg.JWTSecret,
	}
	adHandlers := &handlers.AdHandlers{
		Ads:       store,
		Users:     store,
		JWTSecret: cfg.JWTSecret,
	}
	r := router.NewRouter(authHandlers, adHandlers, cfg.JWTSecret)
	return &App{Cfg: cfg, Router: r}, nil
}

func (a *App) Run() error {
	log.Printf("Server started on %s", a.Cfg.Port)
	return http.ListenAndServe(a.Cfg.Port, a.Router)
}
