package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestHandleGenerateQr_Success(t *testing.T) {
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
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	var errResp ErrorResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
	if errResp.Error != "URL is required" {
		t.Errorf("unexpected error message: %s", errResp.Error)
	}
}
