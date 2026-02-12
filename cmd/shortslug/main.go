package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/StealthBadger747/ShortSlug/internal/bot"
	"github.com/StealthBadger747/ShortSlug/internal/server"
	"github.com/StealthBadger747/ShortSlug/internal/store/sqlite"
)

func main() {
	var (
		defaultPort     = envOrDefault("SERVER_PORT", "8080")
		defaultFrontend = envOrDefault("FRONTEND_DIR", "")
		defaultDB       = envOrDefault("DATABASE_PATH", "")
	)

	port := flag.String("port", defaultPort, "server port")
	frontendDir := flag.String("frontend", defaultFrontend, "path to frontend assets")
	dbPath := flag.String("db", defaultDB, "path to sqlite database file")
	flag.Parse()

	if *frontendDir == "" {
		if dirExists("static") {
			*frontendDir = "static"
		} else {
			log.Fatalf("frontend directory not set; use FRONTEND_DIR or -frontend")
		}
	}

	if *dbPath == "" {
		*dbPath = "shortslug.db"
	}

	absFrontend, err := filepath.Abs(*frontendDir)
	if err != nil {
		log.Fatalf("failed to resolve frontend directory: %v", err)
	}

	store, err := sqlite.Open(*dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer store.Close()

	capVerifier := &bot.CapVerifier{
		SiteVerifyURL: envOrDefault("CAP_SITEVERIFY_URL", ""),
		Secret:        envOrDefault("CAP_SECRET", ""),
	}
	capAPIEndpoint := envOrDefault("CAP_API_ENDPOINT", "")
	publicBaseURL := envOrDefault("PUBLIC_BASE_URL", "")
	password := envOrDefault("SHORTEN_PASSWORD", "")
	brandName := envOrDefault("BRAND_NAME", "ShortSlug")
	analyticsPassword := envOrDefault("ANALYTICS_PASSWORD", "")

	srv := &http.Server{
		Addr:              ":" + *port,
		Handler:           server.New(absFrontend, store, capVerifier, capAPIEndpoint, publicBaseURL, password, brandName, analyticsPassword),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
