package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/grpcreflect"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"reverie.jp/reverie/internal/config"
	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/logger"
)

func getDialOptions(cfg *config.Config) []grpc.DialOption {
	if cfg.Env == config.EnvDevelopment {
		return []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	}
	// TODO: 本番 TLS 証明書の設定
	return []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	}
}

func Run() error {
	cfg := config.New()

	if err := cfg.LoadFromEnv(); err != nil {
		return fmt.Errorf("failed to load config from env: %w", err)
	}

	logger.Init(cfg)

	ctx := context.Background()

	poolCfg, err := pgxpool.ParseConfig(cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("failed to parse database DSN: %w", err)
	}

	poolCfg.MaxConns = cfg.Database.MaxConns
	poolCfg.MinConns = cfg.Database.MinConns
	poolCfg.MaxConnLifetime = cfg.Database.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.Database.MaxConnIdleTime

	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}
	defer db.Close()

	jwtManager := jwt.NewManager(cfg.Auth.JWTSecretKey, cfg.Auth.AccessExpiration, cfg.Auth.RefreshExpiration)

	services := initServices(cfg, db, jwtManager)

	// initialize server mux

	mux := http.NewServeMux()
	gwMux := runtime.NewServeMux()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	for _, service := range services {
		service.RegisterConnectHandler(mux)
		if err := service.RegisterGatewayHandler(ctx, gwMux, addr, getDialOptions(cfg)); err != nil {
			return fmt.Errorf("failed to register gateway handler for service %s: %w", service.Name, err)
		}
	}

	// enable grpc reflection for debugging in development environment
	if cfg.Env == config.EnvDevelopment {
		serviceNames := make([]string, 0, len(services))
		for _, service := range services {
			serviceNames = append(serviceNames, service.Name)
		}
		mux.Handle(grpcreflect.NewHandlerV1(grpcreflect.NewStaticReflector(serviceNames...)))
	}

	mux.Handle("/", gwMux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"Connect-Protocol-Version",
		},
	})

	handler := c.Handler(h2c.NewHandler(mux, &http2.Server{}))

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info("api server is running",
		slog.String("addr", server.Addr),
		slog.String("env", string(cfg.Env)),
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	slog.Info("shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("server exited properly")

	return nil
}
