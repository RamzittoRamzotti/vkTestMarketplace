package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"vkTestMarketplace/internal/http-server/handlers"
	"vkTestMarketplace/internal/http-server/middleware/auth"
	"vkTestMarketplace/internal/storage/sqlite"
)

func main() {
	_ = godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("jwt secret not set")
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./storage/storage.db"
	}

	store, err := sqlite.New(dbPath)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	h := &handlers.AuthHandlers{
		Users:     store,
		JWTSecret: jwtSecret,
	}

	adHandlers := &handlers.AdHandlers{
		Ads:       store,
		Users:     store,
		JWTSecret: jwtSecret,
	}
	adHandlerWithMiddleware := &auth.AdHandler{
		Handlers: adHandlers,
	}
	r := chi.NewRouter()
	r.Use(auth.OptionalAuthMiddleware(jwtSecret))

	r.Post("/register", h.RegisterHandler)
	r.Post("/login", h.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(adHandlerWithMiddleware.AuthMiddleware)
		r.Post("/ads", adHandlers.CreateAdHandler)
	})

	r.Get("/ads", adHandlers.ListAdsHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
