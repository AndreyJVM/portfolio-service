package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/signintech/gopdf"
	qrcode "github.com/skip2/go-qrcode"
)

//go:embed fonts/DejaVuSans.ttf
var fontRegularBytes []byte

//go:embed fonts/DejaVuSans-Bold.ttf
var fontBoldBytes []byte

func main() {
	// Native Docker healthcheck mode: /app/pdf-service healthcheck
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		runHealthcheckCli()
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /api/pdf/health", handleHealth)
	mux.HandleFunc("GET /api/pdf/resume", handleResumePdf)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(corsMiddleware(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("🚀 PDF Service (Go) starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func runHealthcheckCli() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
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
		"service": "pdf-service",
	})
}

func handleResumePdf(w http.ResponseWriter, r *http.Request) {
	pdfBytes, err := generateResumePdf()
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to generate resume PDF",
		})
		return
	}

	disposition := "inline"
	if r.URL.Query().Get("download") == "true" {
		disposition = "attachment"
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=\"Andrey_Vorobev_QA_DevOps_Resume.pdf\"", disposition))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfBytes)
}

func generateResumePdf() ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err := pdf.AddTTFFontData("sans", fontRegularBytes); err != nil {
		return nil, fmt.Errorf("failed to load regular font: %w", err)
	}
	if err := pdf.AddTTFFontData("sans-bold", fontBoldBytes); err != nil {
		return nil, fmt.Errorf("failed to load bold font: %w", err)
	}

	const leftMargin = 40.0
	const pageWidth = 595.28 // A4 width in pt
	const contentWidth = pageWidth - (leftMargin * 2)

	// ==================== HEADER ====================\n
	pdf.SetY(35)
	pdf.SetX(leftMargin)
	_ = pdf.SetFont("sans-bold", "", 20)
	pdf.SetTextColor(15, 23, 42) // Slate 900
	_ = pdf.Cell(nil, "АНДРЕЙ ВОРОБЬЕВ")

	pdf.SetY(60)
	pdf.SetX(leftMargin)
	_ = pdf.SetFont("sans-bold", "", 12)
	pdf.SetTextColor(37, 99, 235) // Blue 600
	_ = pdf.Cell(nil, "DevOps & Linux Infrastructure Engineer")

	pdf.SetY(78)
	pdf.SetX(leftMargin)
	_ = pdf.SetFont("sans", "", 9)
	pdf.SetTextColor(71, 85, 105) // Slate 600
	_ = pdf.Cell(nil, "🌐 Сайт: https://vorobevaqa.ru   |   📱 Telegram: @AndreyAQA   |   💻 GitHub: github.com/AndreyJVM")

	// Embedded QR Code pointing to site (top right)
	qrBytes, err := qrcode.Encode("https://vorobevaqa.ru", qrcode.Medium, 100)
	if err == nil {
		imgHolder, imgErr := gopdf.ImageHolderByBytes(qrBytes)
		if imgErr == nil {
			qrX := pageWidth - leftMargin - 52
			_ = pdf.ImageByHolder(imgHolder, qrX, 32, &gopdf.Rect{W: 52, H: 52})
			pdf.SetY(86)
			pdf.SetX(qrX - 2)
			_ = pdf.SetFont("sans", "", 7)
			pdf.SetTextColor(100, 116, 139)
			_ = pdf.Cell(nil, "vorobevaqa.ru")
		}
	}

	// Divider line
	pdf.SetY(96)
	pdf.SetLineWidth(1)
	pdf.SetStrokeColor(226, 232, 240) // Slate 200
	pdf.Line(leftMargin, 96, leftMargin+contentWidth, 96)

	currentY := 108.0

	// Helper: Section title with stylish accent line
	drawSection := func(title string) {
		pdf.SetY(currentY)
		pdf.SetX(leftMargin)
		_ = pdf.SetFont("sans-bold", "", 11)
		pdf.SetTextColor(30, 41, 59) // Slate 800
		_ = pdf.Cell(nil, title)

		pdf.SetLineWidth(1.5)
		pdf.SetStrokeColor(37, 99, 235) // Blue 600
		pdf.Line(leftMargin, currentY+15, leftMargin+28, currentY+15)

		currentY += 22
	}

	// ==================== SECTION 1: SUMMARY ====================
	drawSection("О СЕБЕ / PROFESSIONAL SUMMARY")
	_ = pdf.SetFont("sans", "", 9)
	pdf.SetTextColor(51, 65, 85)

	summaryLines := []string{
		"DevOps & Linux Infrastructure Engineer с практическим опытом системного администрирования серверов Linux,",
		"контейнеризации в Docker, настройки реверс-прокси Nginx, автоматизации серверной инфраструктуры через Ansible",
		"и построения CI/CD пайплайнов (GitHub Actions, GitLab CI). Инженерный бэкграунд в экосистеме Java",
		"позволяет глубоко понимать процесс разработки и выстраивать надежные, стабильные процессы релизов.",
	}
	for _, line := range summaryLines {
		pdf.SetY(currentY)
		pdf.SetX(leftMargin)
		_ = pdf.Cell(nil, line)
		currentY += 13
	}
	currentY += 8

	// ==================== SECTION 2: CORE SKILLS ====================
	drawSection("КЛЮЧЕВЫЕ НАВЫКИ И ТЕХНОЛОГИЧЕСКИЙ СТЕК")

	skills := []struct {
		Category string
		Items    string
	}{
		{"Linux & OS:", "Ubuntu / Debian, Bash/Shell, systemd, cron, SSH по ключам, UFW firewall, Fail2ban, cgroups"},
		{"Containers & Proxy:", "Docker, Docker Compose, Docker Swarm, Nginx (Reverse Proxy, SSL, Gzip), Docker Hub"},
		{"CI/CD & IaC:", "GitHub Actions, GitLab CI, Ansible (playbooks/roles), Blue-Green деплой, авто-откаты"},
		{"Backend & Scripting:", "Bash, Java (Spring Boot), Go (утилиты/микросервисы), PostgreSQL, REST API, Git"},
		{"QA & Reliability:", "Healthcheck-мониторинг, автоматизированное тестирование, Allure 2, ротация логов"},
	}

	for _, s := range skills {
		pdf.SetY(currentY)
		pdf.SetX(leftMargin)
		_ = pdf.SetFont("sans-bold", "", 8.5)
		pdf.SetTextColor(37, 99, 235)
		_ = pdf.Cell(nil, "• "+s.Category)

		pdf.SetX(leftMargin + 115)
		_ = pdf.SetFont("sans", "", 8.5)
		pdf.SetTextColor(51, 65, 85)
		_ = pdf.Cell(nil, s.Items)
		currentY += 13
	}
	currentY += 8

	// ==================== SECTION 3: WORK EXPERIENCE ====================
	drawSection("ОПЫТ РАБОТЫ И КЛЮЧЕВЫЕ ДОСТИЖЕНИЯ")

	pdf.SetY(currentY)
	pdf.SetX(leftMargin)
	_ = pdf.SetFont("sans-bold", "", 9.5)
	pdf.SetTextColor(15, 23, 42)
	_ = pdf.Cell(nil, "DevOps & Infrastructure Engineer")

	pdf.SetX(pageWidth - leftMargin - 100)
	_ = pdf.SetFont("sans", "", 8.5)
	pdf.SetTextColor(100, 116, 139)
	_ = pdf.Cell(nil, "2022 — Настоящее время")
	currentY += 14

	expBullets := []string{
		"• Настройка и сопровождение серверов на базе Linux (Ubuntu): сетевая безопасность, UFW, управление пользователями и SSH.",
		"• Контейнеризация сервисов в Docker, настройка изоляции, квот CPU/RAM (cgroups limits) и in-binary healthchecks.",
		"• Построение пайплайнов CI/CD в GitHub Actions и GitLab CI с авто-публикацией образов в Docker Hub Registry.",
		"• Разработка отказоустойчивых скриптов развертывания с автоматическим откатом (Rollback) при сбое проверок здоровья.",
		"• Написание переиспользуемых Ansible-ролей для автоматизации конфигурации Linux-хостов и инфраструктурного окружения.",
	}

	for _, b := range expBullets {
		pdf.SetY(currentY)
		pdf.SetX(leftMargin + 5)
		_ = pdf.SetFont("sans", "", 8.5)
		pdf.SetTextColor(51, 65, 85)
		_ = pdf.Cell(nil, b)
		currentY += 13
	}
	currentY += 8

	// ==================== SECTION 4: PROJECTS ====================
	drawSection("ПРОЕКТЫ / FEATURED PROJECTS")

	pdf.SetY(currentY)
	pdf.SetX(leftMargin)
	_ = pdf.SetFont("sans-bold", "", 9.5)
	pdf.SetTextColor(15, 23, 42)
	_ = pdf.Cell(nil, "Portfolio & Infrastructure Platform (https://vorobevaqa.ru)")
	currentY += 14

	projBullets := []string{
		"• Linux VDS хост: защита UFW/Fail2ban, системный тюнинг, изоляция процессов, управление ресурсами и ротация логов.",
		"• Nginx API Gateway: SSL Let's Encrypt, HTTP/2, обратное проксирование микросервисов и кастомные fallback-страницы.",
		"• Контейнерный стек Docker Compose: микросервисы (Java + легковесные Go-утилиты в scratch), healthchecks и квоты памяти.",
		"• Автоматизация CI/CD: пайплайн GitHub Actions, Blue-Green ротация контейнеров без даунтайма с авто-откатом.",
	}
	for _, b := range projBullets {
		pdf.SetY(currentY)
		pdf.SetX(leftMargin + 5)
		_ = pdf.SetFont("sans", "", 8.5)
		pdf.SetTextColor(51, 65, 85)
		_ = pdf.Cell(nil, b)
		currentY += 13
	}
	currentY += 8

	// ==================== SECTION 5: EDUCATION ====================
	drawSection("ОБРАЗОВАНИЕ И КВАЛИФИКАЦИЯ")

	eduItems := []string{
		"• Университет «Синергия» — Бакалавр (09.03.02 Информационные системы и технологии).",
		"• Профильная сертификация: Слёрм (Администрирование Linux), PurpleSchool (Docker + Ansible), Stepik (GitLab CI/CD).",
	}
	for _, e := range eduItems {
		pdf.SetY(currentY)
		pdf.SetX(leftMargin + 5)
		_ = pdf.SetFont("sans", "", 8.5)
		pdf.SetTextColor(51, 65, 85)
		_ = pdf.Cell(nil, e)
		currentY += 13
	}

	// ==================== FOOTER ====================
	pdf.SetY(780)
	pdf.SetLineWidth(0.5)
	pdf.SetStrokeColor(226, 232, 240)
	pdf.Line(leftMargin, 780, leftMargin+contentWidth, 780)

	pdf.SetY(786)
	pdf.SetX(leftMargin)
	_ = pdf.SetFont("sans", "", 7.5)
	pdf.SetTextColor(148, 163, 184)
	footerText := fmt.Sprintf("Сгенерировано в реальном времени микросервисом pdf-service • https://vorobevaqa.ru • %s", time.Now().Format("02.01.2006"))
	_ = pdf.Cell(nil, footerText)

	pdf.SetX(pageWidth - leftMargin - 60)
	_ = pdf.Cell(nil, "Стр. 1 из 1")

	return pdf.GetBytesPdfReturnErr()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
