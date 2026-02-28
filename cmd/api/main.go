package main

import (
	"context"
	"fmt"
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
// @description Sistem pemantauan SNMP Multi-OLT untuk perangkat ZTE
// @host localhost:8080
// @BasePath /
func main() {
	// 1. Inisialisasi Logger untuk mencetak log ke konsol
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// 2. Memuat konfigurasi dari file atau environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// 3. Menyiapkan koneksi Redis jika dikonfigurasi
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

	// 4. Inisialisasi Service (Logika Bisnis) dan Handler (Pengelola HTTP)
	onuService := service.NewONUService(cfg, redisClient)
	oltService := service.NewOLTService(cfg)
	onuHandler := handler.NewONUHandler(onuService)
	oltHandler := handler.NewOLTHandler(oltService)
	queryHandler := handler.NewQueryHandler()
	cliHandler := handler.NewCLIHandler()

	// 5. Setup Router menggunakan Chi
	router := setupRouter(oltHandler, onuHandler, queryHandler, cliHandler)

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

	log.Info().Str("addr", cfg.Server.Addr()).Str("swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.Server.Port)).Msg("Starting")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server error")
	}
}

func setupRouter(oltHandler *handler.OLTHandler, onuHandler *handler.ONUHandler, queryHandler *handler.QueryHandler, cliHandler *handler.CLIHandler) http.Handler {
	r := chi.NewRouter()

	// Menambahkan Middlewares (Fungsi yang berjalan sebelum handler utama)
	r.Use(chiMiddleware.RequestID)    // Memberikan ID unik untuk setiap request
	r.Use(chiMiddleware.RealIP)       // Mendapatkan IP asli client
	r.Use(chiMiddleware.Logger)       // Mencatat log setiap request HTTP
	r.Use(chiMiddleware.Recoverer)    // Mencegah aplikasi crash jika ada panic
	r.Use(middleware.DefaultCORS())   // Mengizinkan akses dari domain luar (Cross-Origin Resource Sharing)
	r.Use(middleware.BasicAuth())     // Autentikasi Basic Auth
	// r.Use(middleware.NewRateLimiter(20, time.Minute).Middleware) // Batasan 20 request per menit per IP
	// r.Use(chiMiddleware.Timeout(90 * time.Second)) // Batas waktu request maksimal 90 detik

	// Endpoint Dasar
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

	// Statistik Pool Koneksi SNMP
	r.Get("/stats", queryHandler.PoolStats)

	// Dokumentasi Swagger UI
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	r.Get("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Grup API Versi 1
	r.Route("/api/v1", func(r chi.Router) {
		// Endpoint "Stateless" (Tanpa simpan kredensial)
		r.Post("/query", queryHandler.Query)
		r.Post("/olt-info", queryHandler.OLTInfo)

		// CLI Commands via Telnet
		r.Route("/cli", func(r chi.Router) {
			// System
			r.Post("/system/clock", cliHandler.ShowClock)
			
			// Hardware
			r.Post("/card", cliHandler.ShowCard)
			r.Post("/card/slot", cliHandler.ShowCardBySlot)
			r.Post("/rack", cliHandler.ShowRack)
			r.Post("/shelf", cliHandler.ShowShelf)
			r.Post("/subcard", cliHandler.ShowSubCard)
			r.Post("/fan", cliHandler.ShowFan)
			
			// GPON Profiles
			r.Post("/gpon/tcont", cliHandler.ShowTcontProfile)
			r.Post("/gpon/onu-type", cliHandler.ShowOnuType)
			r.Post("/gpon/vlan-profile", cliHandler.ShowVlanProfile)
			r.Post("/gpon/ip-profile", cliHandler.ShowIPProfile)
			r.Post("/gpon/sip-profile", cliHandler.ShowSIPProfile)
			r.Post("/gpon/mgc-profile", cliHandler.ShowMGCProfile)
			
			// Line & Remote Profiles
			r.Post("/profile/line/list", cliHandler.ShowLineProfileList)
			r.Post("/profile/line", cliHandler.ShowLineProfile)
			r.Post("/profile/remote/list", cliHandler.ShowRemoteProfileList)
			r.Post("/profile/remote", cliHandler.ShowRemoteProfile)
			
			// GPON ONU
			r.Post("/onu/state", cliHandler.ShowONUState)
			r.Post("/onu/uncfg", cliHandler.ShowONUUncfg)
			r.Post("/onu/config", cliHandler.ShowONUConfig)
			r.Post("/onu/running", cliHandler.ShowONURunning)
			r.Post("/onu/detail", cliHandler.ShowONUDetail)
			r.Post("/onu/distance", cliHandler.ShowONUDistance)
			r.Post("/onu/traffic", cliHandler.ShowONUTraffic)
			r.Post("/onu/optical", cliHandler.ShowONUOptical)
			
			// VLAN
			r.Post("/vlan/list", cliHandler.ShowVLANList)
			r.Post("/vlan/id", cliHandler.ShowVLANByID)
			
			// Interface
			r.Post("/interface", cliHandler.ShowInterface)
			r.Post("/interface/detail", cliHandler.ShowInterfaceByType)
			r.Post("/interface/mng", cliHandler.ShowMgmtInterface)
			
			// Service Port
			r.Post("/service-port", cliHandler.ShowServicePort)
			
			// IGMP
			r.Post("/igmp", cliHandler.ShowIGMP)
			r.Post("/igmp/mvlan", cliHandler.ShowIGMPMVlan)
			r.Post("/igmp/mvlan/id", cliHandler.ShowIGMPMVlanByID)
			r.Post("/igmp/group", cliHandler.ShowIGMPGroup)
			
			// Users
			r.Post("/user/list", cliHandler.ShowUsers)
			r.Post("/user/online", cliHandler.ShowOnlineUsers)
			
			// WRITE Operations (Provisioning)
			r.Post("/onu/auth", cliHandler.AuthenticateONU)
			r.Post("/onu/delete", cliHandler.DeleteONU)
		})

		// Pengelolaan Data OLT (CRUD) + Operasi ONU
		r.Route("/olts", func(r chi.Router) {
			r.Get("/", oltHandler.List)
			r.Post("/", oltHandler.Create)
			
			// Operasi untuk satu OLT (Get/Update/Delete + ONU operations)
			r.Route("/{olt_id}", func(r chi.Router) {
				// CRUD OLT
				r.Get("/", oltHandler.Get)
				r.Put("/", oltHandler.Update)
				r.Delete("/", oltHandler.Delete)
				
				// ONU Operations
				r.Route("/board/{board_id}/pon/{pon_id}", func(r chi.Router) {
					r.Get("/", onuHandler.List)              // List ONU di satu port PON
					r.Delete("/cache", onuHandler.ClearCache) // Bersihkan cache
					r.Get("/empty", onuHandler.EmptySlots)    // Cek slot kosong
					r.Get("/onu/{onu_id}", onuHandler.Detail) // Detail ONU spesifik
				})
			})
		})
	})

	return r
}
