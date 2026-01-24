package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"nyagoPing/internal/domain/model"
)

func TestFileASCIIArtRepository_Save_Load(t *testing.T) {
	repo := NewFileASCIIArtRepository()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_art.txt")

	lines := []string{
		"  ███╗   ██╗███████╗███╗   ███╗██╗   ██╗",
		"  ████╗  ██║██╔════╝████╗ ████║██║   ██║",
		"  ██╔██╗ ██║█████╗  ██╔████╔██║██║   ██║",
		"  ██║╚██╗██║██╔══╝  ██║╚██╔╝██║██║   ██║",
		"  ██║ ╚████║███████╗██║ ╚═╝ ██║╚██████╔╝",
	}
	art, err := model.NewASCIIArt(lines)
	if err != nil {
		t.Fatalf("NewASCIIArt() error = %v", err)
	}

	if err := repo.Save(testFile, art); err != nil {
		t.Errorf("Save() error = %v", err)
	}

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Save() ファイルが作成されていません")
	}

	loadedArt, err := repo.Load(testFile)
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}

	if loadedArt.LineCount() != art.LineCount() {
		t.Errorf("Load() LineCount = %v, want %v", loadedArt.LineCount(), art.LineCount())
	}

	for i, line := range art.Lines() {
		if loadedArt.GetLine(i) != line {
			t.Errorf("Load() line[%d] = %v, want %v", i, loadedArt.GetLine(i), line)
		}
	}
}

func TestFileASCIIArtRepository_Load_NotFound(t *testing.T) {
	repo := NewFileASCIIArtRepository()

	_, err := repo.Load("non_existent_file.txt")
	if err == nil {
		t.Error("Load() 存在しないファイルでエラーが発生しませんでした")
	}
}
