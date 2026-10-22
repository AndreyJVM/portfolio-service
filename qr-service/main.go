package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /api/qr", handleGenerateQr)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(corsMiddleware(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("🚀 QR Code Service (Go) starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func runHealthcheckCli() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%s/health", port))
	if err != nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			log.Printf("Healthcheck failed: %v", err)
		} else {
			log.Printf("Healthcheck failed with status: %d", resp.StatusCode)
		}
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
	w.Header().Set("Content-Type", "application/json")

	var req QrRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "URL is required")
		return
	}

	trimmedURL := strings.TrimSpace(req.URL)
	if trimmedURL == "" {
		sendError(w, http.StatusBadRequest, "URL is required")
		return
	}

	if !urlProtocolRegex.MatchString(trimmedURL) {
		sendError(w, http.StatusBadRequest, "URL must start with http:// or https://")
		return
	}

	pngBytes, err := qrcode.Encode(trimmedURL, qrcode.Medium, 300)
	if err != nil {
		log.Printf("Failed to generate QR for url %s: %v", trimmedURL, err)
		sendError(w, http.StatusInternalServerError, "Failed to generate QR code")
		return
	}

	base64QR := base64.StdEncoding.EncodeToString(pngBytes)

	resp := QrResponse{
		URL:    trimmedURL,
		QRCode: base64QR,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
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
