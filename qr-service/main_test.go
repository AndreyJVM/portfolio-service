package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 from router for /health, got %d", rr.Code)
	}
}

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

func TestHandleGenerateQr_PayloadTooLarge(t *testing.T) {
	// Создаем тело запроса больше 64 КБ
	largeURL := "https://example.com/" + strings.Repeat("a", 70*1024)
	reqBody, _ := json.Marshal(map[string]string{"url": largeURL})
	req := httptest.NewRequest(http.MethodPost, "/api/qr", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handleGenerateQr(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for oversized payload, got %d", rr.Code)
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

func TestHandleHealth_Detailed(t *testing.T) {
	// 1. Виртуальный GET запрос к /health
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	// 2. Вызов хендлера
	handleHealth(rr, req)

	// 3. Проверка HTTP статус-кода
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// 4. Проверка заголовка Content-Type
	contentType := rr.Header().Get("Content-Type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		t.Errorf("expected Content-Type '%s', got '%s'", expectedContentType, contentType)
	}

	// 5. Парсинг JSON-ответа
	var responseBody struct {
		Service string `json:"service"`
		Status  string `json:"status"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	// 6. Проверка полей распарсенного JSON
	if responseBody.Status != "UP" {
		t.Errorf("expected status 'UP', got '%s'", responseBody.Status)
	}
	if responseBody.Service != "qr-service" {
		t.Errorf("expected service 'qr-service', got '%s'", responseBody.Service)
	}
}

func TestHandleGenerateQr_PostSuccess_Detailed(t *testing.T) {
	// Список различных валидных URL, которые должен поддерживать генератор
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "primary domain HTTPS",
			url:  "https://vorobevaqa.ru",
		},
		{
			name: "HTTP with subpath and query parameters",
			url:  "http://example.com/courses/docker?lesson=1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Формируем тело запроса безопасно через json.Marshal
			reqDto := map[string]string{"url": tc.url}
			jsonBytes, err := json.Marshal(reqDto)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}

			// 2. Создаем HTTP POST запрос
			req := httptest.NewRequest(http.MethodPost, "/api/qr", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handleGenerateQr(rr, req)

			// 3. Проверяем HTTP статус 200 OK
			if rr.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
			}

			// 4. Парсим JSON ответ
			var resp struct {
				URL    string `json:"url"`
				QRCode string `json:"qrCode"`
			}

			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode JSON response: %v", err)
			}

			// 5. Проверяем, что url в ответе совпадает с переданным tc.url
			if resp.URL != tc.url {
				t.Errorf("expected url '%s', got '%s'", tc.url, resp.URL)
			}

			// 6. Проверяем, что картинка сгенерирована
			if len(resp.QRCode) == 0 {
				t.Errorf("expected non-empty qrCode for url '%s'", tc.url)
			}
		})
	}
}
