package main

import (
	"log"
	"net/http"

	"github.com/elghazx/go-shorturl/internal/adapters/cache"
	"github.com/elghazx/go-shorturl/internal/adapters/handlers"
	"github.com/elghazx/go-shorturl/internal/adapters/repositories"
	"github.com/elghazx/go-shorturl/internal/config"
	"github.com/elghazx/go-shorturl/internal/core/services"
)

func main() {
	// config
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Failed to initialize config:", err)
	}

	// storage
	defer cfg.DB.Close()
	repo := repositories.NewPostgresRepository(cfg.DB)
	cacheService := cache.NewRedisCache(cfg.Redis)

	// inject service
	urlService := services.NewURLService(repo, cacheService)

	// handlers
	httpHandler := handlers.NewHTTPHandler(urlService)

	// routes
	http.HandleFunc("/", httpHandler.HandleHome)
	http.HandleFunc("/shorten", httpHandler.HandleShorten)
	http.HandleFunc("/stats", httpHandler.HandleStats)
	http.HandleFunc("/api/stats", httpHandler.HandleAPIStats)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
