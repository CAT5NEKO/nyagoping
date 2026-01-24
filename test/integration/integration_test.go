package integration

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"nyagoPing/internal/application/usecase"
	"nyagoPing/internal/domain/model"
	"nyagoPing/internal/domain/service"
	"nyagoPing/internal/infrastructure/persistence"
)

func TestGenerateASCIIArtUseCase_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "generated_art.txt")

	testImagePath := filepath.Join(tmpDir, "test_image.png")
	img := image.NewGray(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.Gray{Y: uint8((x + y) % 256)})
		}
	}
	f, _ := os.Create(testImagePath)
	png.Encode(f, img)
	f.Close()

	repo := persistence.NewFileASCIIArtRepository()
	generator := service.NewASCIIArtGenerator()
	uc := usecase.NewGenerateASCIIArtUseCase(repo, generator)

	input := &usecase.GenerateInput{
		ImagePath:  testImagePath,
		OutputPath: outputPath,
		Width:      40,
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(output.Arts) == 0 {
		t.Error("アスキーアートが生成されませんでした")
	}

	if output.Arts[0].LineCount() == 0 {
		t.Error("生成されたアスキーアートが空です")
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("生成されたファイルが存在しません")
	}

	loadedArt, err := repo.Load(outputPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loadedArt.LineCount() != output.Arts[0].LineCount() {
		t.Errorf("保存されたアートの行数が一致しません: got %v, want %v", loadedArt.LineCount(), output.Arts[0].LineCount())
	}
}

func TestPingUseCase_AutoCount_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	artPath := filepath.Join(tmpDir, "test_art.txt")

	lines := []string{
		"line1",
		"line2",
		"line3",
	}
	art, _ := model.NewASCIIArt(lines)

	repo := persistence.NewFileASCIIArtRepository()
	if err := repo.Save(artPath, art); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	generator := service.NewASCIIArtGenerator()

	expectedCount := generator.CalculateOptimalCount(art)
	if expectedCount != 3 {
		t.Errorf("CalculateOptimalCount() = %v, want 3", expectedCount)
	}
}
