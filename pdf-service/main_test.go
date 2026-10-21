package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected application/json content type, got %s", contentType)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode json body: %v", err)
	}

	if body["status"] != "UP" {
		t.Errorf("expected status 'UP', got '%s'", body["status"])
	}
	if body["service"] != "pdf-service" {
		t.Errorf("expected service 'pdf-service', got '%s'", body["service"])
	}
}

func TestResumePdfHandler_Inline(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/pdf/resume", nil)
	w := httptest.NewRecorder()

	handleResumePdf(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/pdf" {
		t.Errorf("expected 'application/pdf', got '%s'", contentType)
	}

	disposition := w.Header().Get("Content-Disposition")
	if !strings.Contains(disposition, "inline") || !strings.Contains(disposition, "Andrey_Vorobev_QA_DevOps_Resume.pdf") {
		t.Errorf("unexpected Content-Disposition: %s", disposition)
	}

	cacheControl := w.Header().Get("Cache-Control")
	if cacheControl != "public, max-age=3600" {
		t.Errorf("unexpected Cache-Control: %s", cacheControl)
	}

	data := w.Body.Bytes()
	if len(data) < 1000 {
		t.Fatalf("PDF payload too small: %d bytes", len(data))
	}

	// Validate PDF magic bytes: %PDF-
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		t.Errorf("expected PDF magic bytes '%%PDF-', got: %s", string(data[:10]))
	}
}

func TestResumePdfHandler_Attachment(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/pdf/resume?download=true", nil)
	w := httptest.NewRecorder()

	handleResumePdf(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	disposition := w.Header().Get("Content-Disposition")
	if !strings.Contains(disposition, "attachment") {
		t.Errorf("expected attachment in Content-Disposition, got '%s'", disposition)
	}
}

func TestGenerateResumePdf(t *testing.T) {
	pdfBytes, err := generateResumePdf()
	if err != nil {
		t.Fatalf("generateResumePdf() returned error: %v", err)
	}

	if len(pdfBytes) < 5000 {
		t.Errorf("expected generated PDF size > 5000 bytes, got %d", len(pdfBytes))
	}

	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Errorf("PDF does not start with %%PDF- magic bytes")
	}
}

func TestCorsMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := corsMiddleware(nextHandler)

	// Test OPTIONS preflight
	reqOptions := httptest.NewRequest(http.MethodOptions, "/api/pdf/resume", nil)
	wOptions := httptest.NewRecorder()
	handler.ServeHTTP(wOptions, reqOptions)

	if wOptions.Code != http.StatusOK {
		t.Errorf("expected 200 on OPTIONS, got %d", wOptions.Code)
	}
	if wOptions.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("missing or invalid Access-Control-Allow-Origin")
	}

	// Test GET request
	reqGet := httptest.NewRequest(http.MethodGet, "/api/pdf/resume", nil)
	wGet := httptest.NewRecorder()
	handler.ServeHTTP(wGet, reqGet)

	if wGet.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("missing Access-Control-Allow-Origin on GET")
	}
}
