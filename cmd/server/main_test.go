package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateDirs(t *testing.T) {
	tempDir := t.TempDir()
	dir1 := filepath.Join(tempDir, "content")
	dir2 := filepath.Join(tempDir, "templates")

	if err := os.MkdirAll(dir1, 0755); err != nil {
		t.Fatalf("Не удалось создать dir1: %v", err)
	}
	if err := os.MkdirAll(dir2, 0755); err != nil {
		t.Fatalf("Не удалось создать dir2: %v", err)
	}

	regularFile := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(regularFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Не удалось создать файл: %v", err)
	}

	t.Run("Valid existing directories", func(t *testing.T) {
		err := validateDirs(dir1, dir2)
		if err != nil {
			t.Errorf("Ожидался успешный результат, получена ошибка: %v", err)
		}
	})

	t.Run("Non-existent directory returns error", func(t *testing.T) {
		nonExistent := filepath.Join(tempDir, "missing-dir")
		err := validateDirs(dir1, nonExistent)
		if err == nil {
			t.Errorf("Ожидалась ошибка для несуществующего пути")
		}
	})

	t.Run("Path is a file not directory returns error", func(t *testing.T) {
		err := validateDirs(regularFile)
		if err == nil {
			t.Errorf("Ожидалась ошибка, так как путь указывает на обычный файл")
		}
	})
}
