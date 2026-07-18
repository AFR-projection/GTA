package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/AFR-projection/GTA/backend/internal/api"
	"github.com/AFR-projection/GTA/backend/internal/auth"
	"github.com/AFR-projection/GTA/backend/internal/config"
	"github.com/AFR-projection/GTA/backend/internal/db"
	"github.com/AFR-projection/GTA/backend/internal/store"
	mw "github.com/AFR-projection/GTA/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	loadDotEnv()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	migrationsDir := filepath.Join(findBackendRoot(), "migrations")
	if err := db.Migrate(cfg.DatabaseURL, migrationsDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.JWTTTL)
	h := &api.Handler{
		Store:  store.New(pool),
		Tokens: tokens,
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(mw.RateLimit(180))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))
	r.Mount("/", h.Routes())

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("ILO API listening on %s (%s)", cfg.HTTPAddr, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("ILO API stopped")
}

func loadDotEnv() {
	// Prefer repo-root .env, then backend/.env
	candidates := []string{
		filepath.Join("..", ".env"),
		".env",
		filepath.Join("..", "..", ".env"),
	}
	for _, c := range candidates {
		if err := godotenv.Load(c); err == nil {
			log.Printf("loaded env from %s", c)
			return
		}
	}
}

func findBackendRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	// If running from backend/ or backend/cmd/api
	if _, err := os.Stat(filepath.Join(wd, "migrations")); err == nil {
		return wd
	}
	if _, err := os.Stat(filepath.Join(wd, "..", "..", "migrations")); err == nil {
		return filepath.Clean(filepath.Join(wd, "..", ".."))
	}
	if _, err := os.Stat(filepath.Join(wd, "..", "migrations")); err == nil {
		return filepath.Clean(filepath.Join(wd, ".."))
	}
	return wd
}
