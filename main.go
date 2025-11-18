package main

import (
	_ "github.com/GoAdminGroup/go-admin/adapter/gin"                 // web framework adapter
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres" // sql driver
	_ "github.com/GoAdminGroup/themes/sword"                         // ui theme

	"Clinic_backend/cmd/app"
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	_ "Clinic_backend/docs"
)

// @title Clinic Backend API
// @version 1.0
// @description API для управления частной клиникой
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@clinic.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {

	//fmt.Println("Running main...")
	slog.Info("Running main...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Канал для получения ошибок из приложения
	appError := make(chan error, 1)

	go func() {
		if err := app.StartApplication(ctx); err != nil {
			appError <- err
		} else {
			appError <- nil
		}
	}()

	// Ожидаем либо сигнал ОС, либо ошибку из приложения
	select {
	case <-ctx.Done():
		// Получен сигнал от ОС
		slog.Info("Received shutdown signal, shutting down gracefully...")

	case err := <-appError:
		// Получена ошибка из приложения
		if err != nil {
			slog.Error("Application error, shutting down gracefully", "error", err.Error())
		} else {
			slog.Info("Application completed successfully")
		}
	}

	slog.Info("Shutdown completed")
}
