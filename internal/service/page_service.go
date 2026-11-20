package service

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"

	"ru.vorobevaqa/devops-platform/internal/model"
)

// PageMeta defines the structure for page metadata stored in JSON.
type PageMeta struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Keywords      string `json:"keywords"`
	OGTitle       string `json:"ogTitle,omitempty"`
	OGDescription string `json:"ogDescription,omitempty"`
	OGImage       string `json:"ogImage,omitempty"`
	ActiveTab     string `json:"activeTab,omitempty"`
}

// PageService manages page metadata and content loaded from JSON.
type PageService struct {
	pagesFilePath string
	mu            sync.RWMutex
	pages         map[string]PageMeta
}

// defaultPages provides built-in fallback metadata in case the external configuration is unavailable.
var defaultPages = map[string]PageMeta{
	"index": {
		Title:         "Об авторе • Андрей Воробьев | Primary Cloud Lab",
		Description:   "Об авторе платформы Primary Cloud Lab Андрее Воробьеве: инженерный опыт, системное администрирование Linux (Слёрм), Docker, Podman, LXC, CI/CD, репозитории и квалификация.",
		Keywords:      "Андрей Воробьев, об авторе, DevOps, Linux, Слёрм, Docker, Podman, LXC, Ansible, сертификаты, квалификация, инфраструктура",
		OGTitle:       "Об авторе • Андрей Воробьев | Primary Cloud Lab",
		OGDescription: "Создатель Primary Cloud Lab. Архитектура инфраструктуры, open-source репозитории, профессиональные сертификации Слёрм и дипломы.",
		OGImage:       "/images/og-about.png",
		ActiveTab:     "about",
	},
	"mustang": {
		Title:         "На мечту • Ford Mustang Shelby GT500",
		Description:   "Фонд спасения от рутины: краудфандинг на мечту — Ford Mustang Shelby GT500.",
		Keywords:      "На мечту, Ford Mustang, Shelby GT500, краудфандинг, V8",
		OGTitle:       "На мечту • Ford Mustang Shelby GT500",
		OGDescription: "Осталось закрыть единственный критический баг — отсутствие 760 л.с. в гараже.",
		OGImage:       "/images/avatar.jpg",
		ActiveTab:     "mustang",
	},
	"404": {
		Title:         "404 • Страница не найдена",
		Description:   "Запрошенная страница не найдена на платформе Primary Cloud Lab.",
		Keywords:      "404, страница не найдена, devops, linux",
		OGTitle:       "404 • Страница не найдена",
		OGDescription: "Запрошенная страница не найдена на платформе Primary Cloud Lab.",
		OGImage:       "/images/og-main.png",
		ActiveTab:     "",
	},
}

// NewPageService initializes a new PageService and loads metadata from the given JSON file.
func NewPageService(pagesFilePath string) *PageService {
	ps := &PageService{
		pagesFilePath: pagesFilePath,
		pages:         make(map[string]PageMeta),
	}

	if err := ps.loadPages(); err != nil {
		log.Printf("⚠️ Не удалось загрузить метаданные страниц из '%s' (%v), используются встроенные значения", pagesFilePath, err)
		ps.useDefaults()
	}

	return ps
}

func (s *PageService) useDefaults() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range defaultPages {
		s.pages[k] = v
	}
}

func (s *PageService) loadPages() error {
	if s.pagesFilePath == "" {
		s.useDefaults()
		return nil
	}

	data, err := os.ReadFile(s.pagesFilePath)
	if err != nil {
		return err
	}

	var pages map[string]PageMeta
	if err := json.Unmarshal(data, &pages); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.pages = pages
	return nil
}

// GetMetaData builds model.MetaData for a given page key, request path, and application base URL.
func (s *PageService) GetMetaData(pageKey, path, baseURL string) model.MetaData {
	s.mu.RLock()
	meta, exists := s.pages[pageKey]
	s.mu.RUnlock()

	if !exists {
		// Fallback to defaultPages if not in custom config
		if def, ok := defaultPages[pageKey]; ok {
			meta = def
		}
	}

	ogTitle := meta.OGTitle
	if ogTitle == "" {
		ogTitle = meta.Title
	}

	ogDesc := meta.OGDescription
	if ogDesc == "" {
		ogDesc = meta.Description
	}

	ogImage := meta.OGImage
	if ogImage != "" && strings.HasPrefix(ogImage, "/") {
		ogImage = baseURL + ogImage
	}

	fullURL := baseURL + path

	return model.MetaData{
		Title:         meta.Title,
		Description:   meta.Description,
		Keywords:      meta.Keywords,
		OGTitle:       ogTitle,
		OGDescription: ogDesc,
		OGImage:       ogImage,
		OGURL:         fullURL,
		CanonicalURL:  fullURL,
		ActiveTab:     meta.ActiveTab,
	}
}
