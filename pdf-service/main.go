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

	// ==================== HEADER ====================
	pdf.SetY(35)
	pdf.SetX(leftMargin)
	_ = pdf.SetFont("sans-bold", "", 20)
	pdf.SetTextColor(15, 23, 42) // Slate 900
	_ = pdf.Cell(nil, "АНДРЕЙ ВОРОБЬЕВ")

	pdf.SetY(60)
	pdf.SetX(leftMargin)
	_ = pdf.SetFont("sans-bold", "", 12)
	pdf.SetTextColor(37, 99, 235) // Blue 600
	_ = pdf.Cell(nil, "Lead QA Automation & DevOps / Backend Engineer")

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
		"Инженер по автоматизации тестирования и DevOps с глубокой экспертизой в экосистеме Java/JVM (Java 21,",
		"Virtual Threads, Spring Boot 3.4), высокопроизводительных микросервисах (Go 1.23, Docker scratch), нагрузочном",
		"тестировании (Apache JMeter 5.5) и автоматизации CI/CD (GitHub Actions). Автор продакшн-архитектуры",
		"сервиса vorobevaqa.ru со 100% покрытием критических бизнес-процессов (125+ automated tests).",
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
		{"Backend & JVM:", "Java 21 (Virtual Threads), Spring Boot 3.4, Go 1.23, Python, REST API, Thymeleaf, Caffeine Cache"},
		{"QA Automation:", "Selenide 7 (Page Object Model), REST Assured, JUnit 5, Allure Report 2, Mockito, AssertJ"},
		{"Load & Performance:", "Apache JMeter 5.5, high-concurrency API plans, HTML Dashboard Reports, Latency SLA <15ms"},
		{"DevOps & Infra:", "Docker, Docker Compose, Nginx (API Gateway, HTTP/2, SSL, Gzip), GitHub Actions, Linux (Ubuntu), Bash"},
		{"Code Quality:", "Spotless (Google Java Style), Checkstyle, Clean Architecture, Zero-downtime Blue/Green Deploy"},
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
	_ = pdf.Cell(nil, "Lead QA Automation / DevOps Engineer")

	pdf.SetX(pageWidth - leftMargin - 100)
	_ = pdf.SetFont("sans", "", 8.5)
	pdf.SetTextColor(100, 116, 139)
	_ = pdf.Cell(nil, "2022 — Настоящее время")
	currentY += 14

	expBullets := []string{
		"• Разработка и внедрение архитектуры автоматизации тестирования с нуля (Java 21, Selenide, REST Assured, Allure).",
		"• Проектирование микросервисного стека с выделенным Go микросервисом (<10MB RAM) и Nginx API Gateway.",
		"• Построение полного цикла CI/CD в GitHub Actions с параллельным тестированием (Maven + Go) и Zero-downtime деплоем.",
		"• Разработка профилей нагрузочного тестирования JMeter, выявление узких мест и оптимизация задержек до <15 мс.",
		"• Обеспечение стабильности и качества кодовой базы: 125+ сквозных тестов, 0 дефектов форматирования (Spotless/Checkstyle).",
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
	_ = pdf.Cell(nil, "Portfolio & Multi-Microservices Platform (https://vorobevaqa.ru)")
	currentY += 14

	projBullets := []string{
		"• High-Performance QR Microservice (Go 1.23, scratch image ~10MB, <15ms response latency).",
		"• Dynamic PDF Generator Microservice (Go 1.23, vector typography, dynamic QR embedding, <20ms generation).",
		"• Java Portfolio Web Application (Spring Boot 3.4, Java 21 Virtual Threads, Caffeine Cache, Bootstrap 5.3).",
		"• Complete QA Test Suite: 125+ automated UI, API and Unit tests, Allure 2 reports, JMeter load dashboards.",
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
		"• Высшее техническое образование в сфере информационных технологий и автоматизированных систем.",
		"• Профильная подготовка: Java Software Engineering, Enterprise AQA Frameworks, Linux Administration & DevOps.",
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
