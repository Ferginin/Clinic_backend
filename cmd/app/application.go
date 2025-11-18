package app

import (
	"Clinic_backend/config"
	"Clinic_backend/internal/router"
	"Clinic_backend/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func StartApplication(ctx context.Context) error {
	cfg := config.GetConfig()

	// Инициализация логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Clinic Backend API", "version", "1.0.0")

	cfg.Client = storage.NewConnection(ctx, cfg)

	// Настройка маршрутов
	r := router.SetupRouter(cfg, cfg.Client)

	addr := fmt.Sprintf("%s:%d", cfg.Env.IpAddress, cfg.Env.API_PORT)
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
