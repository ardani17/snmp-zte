package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardani/snmp-zte/internal/config"
	_ "github.com/ardani/snmp-zte/docs"
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
)

// @title SNMP-ZTE API
// @version 2.1
// @description Multi-OLT SNMP monitoring system for ZTE devices
// @host localhost:8080
// @BasePath /
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	var redisClient *redis.Client
	if cfg.Redis.Host != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			redisClient = nil
		} else {
			defer redisClient.Close()
		}
	}

	onuService := service.NewONUService(cfg, redisClient)
	oltService := service.NewOLTService(cfg)
	onuHandler := handler.NewONUHandler(onuService)
	oltHandler := handler.NewOLTHandler(oltService)
	queryHandler := handler.NewQueryHandler()

	router := setupRouter(oltHandler, onuHandler, queryHandler)

	server := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Info().Msg("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Info().Str("addr", cfg.Server.Addr()).Str("swagger", "http://"+cfg.Server.Addr()+"/swagger/index.html").Msg("Starting")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server error")
	}
}

func setupRouter(oltHandler *handler.OLTHandler, onuHandler *handler.ONUHandler, queryHandler *handler.QueryHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.NewRateLimiter(20, time.Minute).Middleware)
	r.Use(middleware.DefaultCORS())
	r.Use(chiMiddleware.Timeout(90 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"name":    "SNMP-ZTE API",
			"version": "2.1.0",
			"swagger": "/swagger/index.html",
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	})

	r.Get("/stats", queryHandler.PoolStats)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/query", queryHandler.Query)
		r.Post("/olt-info", queryHandler.OLTInfo)

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
