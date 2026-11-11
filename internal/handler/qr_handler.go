package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	"ru.vorobevaqa/devops-platform/internal/model"
)

// QRPage: QR Code Playground страница
func (h *Handler) QRPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Meta model.MetaData
	}{
		Meta: model.MetaData{
			Title:         "QR Code Playground",
			Description:   "Генератор QR-кодов на легковесном нативном Go движке.",
			Keywords:      "QR code generator, golang, qrcode",
			OGTitle:       "QR Code Playground",
			OGDescription: "Интерактивный генератор QR-кодов на Go.",
			OGImage:       "/images/og-qr.png",
			OGURL:         h.baseURL + "/qr",
			ActiveTab:     "qr",
			CanonicalURL:  h.baseURL + "/qr",
		},
	}

	h.renderTemplate(w, "qr.html", data)
}

// GenerateQREndpoint: Native Go QR Code generation (PNG or JSON).
// Защищен от DoS ограничением размера payload (64 КБ).
func (h *Handler) GenerateQREndpoint(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		text = r.URL.Query().Get("url")
	}

	if r.Method == http.MethodPost {
		// Ограничиваем размер входящего тела 64 КБ
		r.Body = http.MaxBytesReader(w, r.Body, 64*1024)

		var body struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil && body.URL != "" {
			text = body.URL
		}
	}

	if text == "" {
		text = h.baseURL
	}

	pngBytes, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Не удалось сгенерировать QR-код", http.StatusInternalServerError)
		return
	}

	acceptsJSON := strings.Contains(r.Header.Get("Accept"), "application/json") ||
		strings.Contains(r.Header.Get("Content-Type"), "application/json") ||
		r.Method == http.MethodPost ||
		r.URL.Query().Get("format") == "json"

	if acceptsJSON {
		base64Str := base64.StdEncoding.EncodeToString(pngBytes)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"qrCode": base64Str,
			"url":    text,
			"size":   256,
			"format": "PNG",
		})
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pngBytes)
}
