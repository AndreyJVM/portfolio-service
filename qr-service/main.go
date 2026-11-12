package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

type QrRequest struct {
	URL string `json:"url"`
}

type QrResponse struct {
	URL    string `json:"url"`
	QRCode string `json:"qrCode"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var urlProtocolRegex = regexp.MustCompile(`^(http|https)://.+$`)

// setupRouter конфигурирует маршруты и middleware микросервиса QR
func setupRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("/api/qr", handleGenerateQr)

	return loggingMiddleware(corsMiddleware(mux))
}

func main() {
	// Native Docker healthcheck mode: /app/qr-service healthcheck
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		runHealthcheckCli()
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      setupRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown по сигналам SIGINT и SIGTERM
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 5*time.Second)
		defer cancel()

		log.Println("🛑 Stopping QR service gracefully...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error shutting down QR service: %v", err)
		}
		serverStopCtx()
	}()

	log.Printf("🚀 QR Code Service (Go) starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	<-serverCtx.Done()
	log.Println("✅ QR service stopped gracefully.")
}

func checkHealth(baseURL string) error {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("healthcheck failed with status: %d", resp.StatusCode)
	}
	return nil
}

func runHealthcheckCli() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := checkHealth(fmt.Sprintf("http://127.0.0.1:%s", port)); err != nil {
		log.Printf("Healthcheck failed: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "UP",
		"service": "qr-service",
	})
}

func handleGenerateQr(w http.ResponseWriter, r *http.Request) {
	var inputURL string

	if r.Method == http.MethodGet {
		inputURL = r.URL.Query().Get("url")
		if inputURL == "" {
			inputURL = r.URL.Query().Get("text")
		}
	} else if r.Method == http.MethodPost {
		// Защита от DoS: ограничиваем размер входящего тела 64 КБ
		r.Body = http.MaxBytesReader(w, r.Body, 64*1024)

		var req QrRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "Invalid JSON body or payload too large")
			return
		}
		inputURL = req.URL
	} else {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	trimmedURL := strings.TrimSpace(inputURL)
	if trimmedURL == "" {
		sendError(w, http.StatusBadRequest, "URL is required")
		return
	}

	if !urlProtocolRegex.MatchString(trimmedURL) {
		sendError(w, http.StatusBadRequest, "URL must start with http:// or https://")
		return
	}

	const maxURLTextLen = 2048
	if len(trimmedURL) > maxURLTextLen {
		sendError(w, http.StatusBadRequest, "URL exceeds maximum length of 2048 characters")
		return
	}

	pngBytes, err := qrcode.Encode(trimmedURL, qrcode.Medium, 300)
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
		sendError(w, http.StatusBadRequest, "Failed to generate QR code: invalid content")
		return
	}

	if r.Method == http.MethodGet && strings.Contains(r.Header.Get("Accept"), "image/") {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pngBytes)
		return
	}

	base64QR := base64.StdEncoding.EncodeToString(pngBytes)

	resp := QrResponse{
		URL:    trimmedURL,
		QRCode: base64QR,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
