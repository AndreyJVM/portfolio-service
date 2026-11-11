package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "UP" {
		t.Errorf("expected status UP, got %s", body["status"])
	}
}

func TestHandleGenerateQr_SuccessPost(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]string{"url": "https://vorobevaqa.ru"})
	req := httptest.NewRequest(http.MethodPost, "/api/qr", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handleGenerateQr(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp QrResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if resp.URL != "https://vorobevaqa.ru" {
		t.Errorf("expected url https://vorobevaqa.ru, got %s", resp.URL)
	}

	if len(resp.QRCode) == 0 {
		t.Error("expected non-empty base64 qrCode")
	}
}

func TestHandleGenerateQr_SuccessGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/qr?url=https://vorobevaqa.ru", nil)
	rr := httptest.NewRecorder()

	handleGenerateQr(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp QrResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if resp.URL != "https://vorobevaqa.ru" {
		t.Errorf("expected url https://vorobevaqa.ru, got %s", resp.URL)
	}

	if len(resp.QRCode) == 0 {
		t.Error("expected non-empty base64 qrCode")
	}
}

func TestHandleGenerateQr_PngImageAccept(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/qr?url=https://vorobevaqa.ru", nil)
	req.Header.Set("Accept", "image/png")
	rr := httptest.NewRecorder()

	handleGenerateQr(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", rr.Header().Get("Content-Type"))
	}
	if rr.Body.Len() == 0 {
		t.Error("expected non-empty png response body")
	}
}

func TestHandleGenerateQr_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/qr", nil)
	rr := httptest.NewRecorder()

	handleGenerateQr(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}

func TestHandleGenerateQr_InvalidProtocol(t *testing.T) {
	invalidURLs := []string{"ftp://test.com", "not-a-url", "javascript:alert(1)"}

	for _, u := range invalidURLs {
		reqBody, _ := json.Marshal(map[string]string{"url": u})
		req := httptest.NewRequest(http.MethodPost, "/api/qr", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()

		handleGenerateQr(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("url %s: expected status 400, got %d", u, rr.Code)
		}

		var errResp ErrorResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
		if errResp.Error != "URL must start with http:// or https://" {
			t.Errorf("unexpected error message: %s", errResp.Error)
		}
	}
}

func TestHandleGenerateQr_BlankUrl(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]string{"url": "   "})
	req := httptest.NewRequest(http.MethodPost, "/api/qr", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handleGenerateQr(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", reqBody)
	}

	var errResp ErrorResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
	if errResp.Error != "URL is required" {
		t.Errorf("unexpected error message: %s", errResp.Error)
	}
}

func TestHandleGenerateQr_InvalidJson(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/qr", strings.NewReader("{invalid-json"))
	rr := httptest.NewRecorder()

	handleGenerateQr(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for malformed json, got %d", rr.Code)
	}
}

func TestMiddlewares(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// 1. Тестируем corsMiddleware (OPTIONS запрос)
	cors := corsMiddleware(testHandler)
	reqOptions := httptest.NewRequest(http.MethodOptions, "/api/qr", nil)
	rrOptions := httptest.NewRecorder()
	cors.ServeHTTP(rrOptions, reqOptions)

	if rrOptions.Code != http.StatusOK {
		t.Errorf("expected 200 for OPTIONS, got %d", rrOptions.Code)
	}
	if rrOptions.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Access-Control-Allow-Origin header missing or invalid")
	}

	// 2. Тестируем loggingMiddleware + corsMiddleware для обычного GET
	handler := loggingMiddleware(cors)
	reqGet := httptest.NewRequest(http.MethodGet, "/health", nil)
	rrGet := httptest.NewRecorder()
	handler.ServeHTTP(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Errorf("expected 200 for GET, got %d", rrGet.Code)
	}
	if rrGet.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got %s", rrGet.Body.String())
	}
}
