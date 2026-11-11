package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"ru.vorobevaqa/devops-platform/internal/config"
	"ru.vorobevaqa/devops-platform/internal/handler"
	"ru.vorobevaqa/devops-platform/internal/service"
)

func main() {
	cfg := config.Load()

	courseService := service.NewCourseService(cfg.ContentDir)
	integrationsService := service.NewIntegrationsService(cfg.GitHubUser, cfg.DockerUser)
	h := handler.NewHandler(courseService, integrationsService, cfg.TemplatesDir, cfg.StaticDir, cfg.AppBaseURL)

	r := chi.NewRouter()

	// Стандартные эффективные middleware Chi
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)
	r.Use(middleware.Compress(5))

	// Healthcheck для Docker Compose (совместим с curl /health или /actuator/health)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"UP","service":"devops-platform","runtime":"golang"}`))
	})
	r.Get("/actuator/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"UP","service":"devops-platform","runtime":"golang"}`))
	})

	// Статические файлы: /css/*, /js/*, /images/*, /pdf/*
	staticFS := http.FileServer(http.Dir(cfg.StaticDir))
	for _, prefix := range []string{"/css", "/js", "/images", "/pdf"} {
		prefixPattern := prefix + "/*"
		r.Handle(prefixPattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=86400")
			staticFS.ServeHTTP(w, r)
		}))
	}

	r.Get("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.StaticDir+"/favicon.svg")
	})
	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.StaticDir+"/robots.txt")
	})
	r.Get("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.StaticDir+"/sitemap.xml")
	})

	// Страницы платформы
	r.Get("/", h.Index)
	r.Get("/courses", h.CoursesList)
	r.Get("/courses/{courseSlug}", h.CourseDetail)
	r.Get("/courses/{courseSlug}/{lessonSlug}", h.Lesson)

	r.Get("/projects", h.Projects)
	r.Get("/education", h.Education)
	r.Get("/about", h.RedirectAbout)
	r.Get("/qr", h.QRPage)
	r.Get("/mustang", h.Mustang)

	// API генерации QR-кодов (Go native engine)
	r.Get("/api/qr", h.GenerateQREndpoint)
	r.Post("/api/qr", h.GenerateQREndpoint)

	// API PDF резюме
	r.Get("/api/pdf/resume", h.ResumePDF)
	r.Get("/api/pdf/health", h.PdfHealth)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 10*time.Second)
		defer cancel()

		log.Println("🛑 Получен сигнал завершения, плавная остановка сервера Go...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Ошибка при остановке сервера: %v\n", err)
		}
		serverStopCtx()
	}()

	log.Printf("🚀 DevOps Platform (Go) успешно запущена на http://0.0.0.0:%s\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка ListenAndServe: %v\n", err)
	}

	<-serverCtx.Done()
	log.Println("✅ Сервер успешно остановлен.")
}
