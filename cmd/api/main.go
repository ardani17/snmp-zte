package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardani/snmp-zte/internal/config"
	"github.com/ardani/snmp-zte/internal/handler"
	"github.com/ardani/snmp-zte/internal/middleware"
	"github.com/ardani/snmp-zte/internal/service"
	"github.com/ardani/snmp-zte/pkg/response"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ardani/snmp-zte/docs" // swagger docs
)

// @title SNMP-ZTE API
// @version 2.1
// @description Multi-OLT SNMP monitoring system for ZTE devices (C320, C300, C600)
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/ardani17/snmp-zte
// @contact.email adifta22@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Initialize Redis client (optional)
	var redisClient *redis.Client
	if cfg.Redis.Host != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			log.Warn().Err(err).Msg("Redis not available, caching disabled")
			redisClient = nil
		} else {
			defer redisClient.Close()
			log.Info().Msg("Redis connected")
		}
	}

	// Initialize services
	onuService := service.NewONUService(cfg, redisClient)
	oltService := service.NewOLTService(cfg)

	// Initialize handlers
	onuHandler := handler.NewONUHandler(onuService)
	oltHandler := handler.NewOLTHandler(oltService)
	queryHandler := handler.NewQueryHandler()

	// Setup router
	router := setupRouter(oltHandler, onuHandler, queryHandler)

	// Start server
	server := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Info().Msg("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Server shutdown error")
		}
		cancel()
	}()

	log.Info().
		Str("addr", cfg.Server.Addr()).
		Str("swagger", "http://"+cfg.Server.Addr()+"/swagger/index.html").
		Int("rate_limit", 20).
		Int("max_concurrent_snmp", 100).
		Msg("Starting server")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server error")
	}

	<-ctx.Done()
	log.Info().Msg("Server stopped")
}

func setupRouter(
	oltHandler *handler.OLTHandler,
	onuHandler *handler.ONUHandler,
	queryHandler *handler.QueryHandler,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(20, time.Minute)
	r.Use(rateLimiter.Middleware)

	// CORS
	r.Use(middleware.DefaultCORS())

	// Timeout
	r.Use(chiMiddleware.Timeout(90 * time.Second))

	// Root endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"name":           "SNMP-ZTE API",
			"version":        "2.1.0",
			"status":         "running",
			"swagger_docs":   "/swagger/index.html",
			"rate_limit":     "20 req/min per IP",
			"max_concurrent": 100,
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// Pool statistics
	r.Get("/stats", queryHandler.PoolStats)

	// Swagger UI - serve at /swagger/*
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Stateless query endpoint
		r.Post("/query", queryHandler.Query)
		r.Post("/olt-info", queryHandler.OLTInfo)

		// Legacy endpoints
		r.Route("/olts", func(r chi.Router) {
			r.Get("/", oltHandler.List)
			r.Post("/", oltHandler.Create)
			r.Get("/{olt_id}", oltHandler.Get)
			r.Put("/{olt_id}", oltHandler.Update)
			r.Delete("/{olt_id}", oltHandler.Delete)
		})

		r.Route("/olts/{olt_id}", func(r chi.Router) {
			r.Route("/board/{board_id}/pon/{pon_id}", func(r chi.Router) {
				r.Get("/", onuHandler.List)
				r.Delete("/cache", onuHandler.ClearCache)
				r.Get("/empty", onuHandler.EmptySlots)
				r.Get("/onu/{onu_id}", onuHandler.Detail)
			})
		})
	})

	return r
}
