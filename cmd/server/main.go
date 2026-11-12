package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ru.vorobevaqa/devops-platform/internal/config"
	"ru.vorobevaqa/devops-platform/internal/handler"
	"ru.vorobevaqa/devops-platform/internal/service"
)

// validateDirs проверяет наличие обязательных директорий ресурсов при старте сервиса.
func validateDirs(dirs ...string) error {
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("каталог '%s' недоступен: %w", dir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("путь '%s' не является каталогом", dir)
		}
	}
	return nil
}

func main() {
	cfg := config.Load()

	// Валидация существования необходимых ресурсов
	if err := validateDirs(cfg.ContentDir, cfg.TemplatesDir, cfg.StaticDir); err != nil {
		log.Printf("⚠️ Предупреждение о путях конфигурации: %v", err)
	}

	// Инициализация сервисного слоя
	courseService := service.NewCourseService(cfg.ContentDir)
	integrationsService := service.NewIntegrationsService(cfg.GitHubUser, cfg.DockerUser)

	// Инициализация контроллеров и роутера
	h := handler.NewHandler(courseService, integrationsService, cfg.TemplatesDir, cfg.StaticDir, cfg.AppBaseURL)
	router := handler.NewRouter(h, cfg.StaticDir)

	server := &http.Server{
		Addr:              net.JoinHostPort("", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown на базе signal.NotifyContext
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("🚀 DevOps Platform (Go) успешно запущена на http://0.0.0.0:%s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		log.Fatalf("❌ Критическая ошибка ListenAndServe: %v\n", err)
	case <-ctx.Done():
		log.Println("🛑 Получен сигнал завершения, плавная остановка сервера Go...")
	}

	// Ожидание завершения активных соединений (grace period 10 секунд)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Ошибка при остановке сервера: %v\n", err)
	} else {
		log.Println("✅ Сервер успешно остановлен.")
	}
}
