package main

import (
	"Clinic_backend/config"
	_ "Clinic_backend/docs"
	"Clinic_backend/internal/router"
	"Clinic_backend/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Clinic Backend API
// @version 1.0
// @description API для управления частной клиникой
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@clinic.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	slog.Info("Running main...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error, 1)

	go func() {
		if err := StartApplication(ctx); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("Received shutdown signal, shutting down gracefully...")

	case err := <-errChan:
		if err != nil {
			slog.Error("Application error, shutting down gracefully", "error", err.Error())
		} else {
			slog.Info("Application completed successfully")
		}
	}

	slog.Info("Shutdown completed")
}

func StartApplication(ctx context.Context) error {
	cfg := config.GetConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Clinic Backend API", "version", "1.0.0")

	cfg.Client = storage.NewConnection(ctx, cfg)
	defer cfg.Client.Close()

	r := router.SetupRouter(cfg, cfg.Client)

	addr := fmt.Sprintf("%s:%d", cfg.Env.IpAddress, cfg.Env.ApiPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		slog.Info("Server started")
		slog.Info("Swagger UI available at", "url", fmt.Sprintf("http://%s/swagger/index.html", addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server %v", slog.String("error", err.Error()))
			panic(err)
		}
	}()

	<-ctx.Done()
	slog.Info("⚫️ Graceful shutdown initiated...")
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("⚫️ Server forced to shutdown", slog.String("error", err.Error()))
		panic(err)
	}
	return nil
}
