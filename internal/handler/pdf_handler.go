package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// ResumePDF отдает PDF-файл резюме
func (h *Handler) ResumePDF(w http.ResponseWriter, r *http.Request) {
	pdfPath := filepath.Join(h.staticDir, "pdf", "resume.pdf")
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		http.Error(w, "Файл резюме не найден", http.StatusNotFound)
		return
	}

	download := r.URL.Query().Get("download") == "true"
	disposition := "inline"
	if download {
		disposition = "attachment"
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=\"Andrey_Vorobev_QA_DevOps_Resume.pdf\"", disposition))
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, pdfPath)
}

// PdfHealth проверяет статус PDF сервиса
func (h *Handler) PdfHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"UP","service":"pdf-service"}`))
}
