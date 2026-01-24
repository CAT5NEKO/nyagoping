package service

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestASCIIArtGenerator_GenerateFromImage(t *testing.T) {
	generator := NewASCIIArtGenerator()

	tmpDir := t.TempDir()
	testImagePath := filepath.Join(tmpDir, "test.png")

	img := image.NewGray(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			brightness := uint8(float64(x) / 100.0 * 255)
			img.Set(x, y, color.Gray{Y: brightness})
		}
	}

	f, err := os.Create(testImagePath)
	if err != nil {
		t.Fatalf("テスト画像の作成エラー: %v", err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatalf("画像のエンコードエラー: %v", err)
	}
	f.Close()

	art, err := generator.GenerateFromImage(testImagePath, 40)
	if err != nil {
		t.Errorf("GenerateFromImage() error = %v", err)
		return
	}

	if art.LineCount() == 0 {
		t.Error("GenerateFromImage() 生成されたアートが空です")
	}
}

func TestASCIIArtGenerator_GenerateFromImagesInDirectory(t *testing.T) {
	generator := NewASCIIArtGenerator()

	tmpDir := t.TempDir()

	for i := 1; i <= 2; i++ {
		testImagePath := filepath.Join(tmpDir, "test"+string(rune('0'+i))+".png")

		img := image.NewGray(image.Rect(0, 0, 50, 50))
		for y := 0; y < 50; y++ {
			for x := 0; x < 50; x++ {
				img.Set(x, y, color.Gray{Y: uint8(i * 50)})
			}
		}

		f, err := os.Create(testImagePath)
		if err != nil {
			t.Fatalf("テスト画像の作成エラー: %v", err)
		}
		png.Encode(f, img)
		f.Close()
	}

	arts, filenames, err := generator.GenerateFromImagesInDirectory(tmpDir, 40)
	if err != nil {
		t.Errorf("GenerateFromImagesInDirectory() error = %v", err)
		return
	}

	if len(arts) != 2 {
		t.Errorf("GenerateFromImagesInDirectory() 生成数 = %d, want 2", len(arts))
	}

	if len(filenames) != 2 {
		t.Errorf("GenerateFromImagesInDirectory() ファイル名数 = %d, want 2", len(filenames))
	}
}

func TestASCIIArtGenerator_CalculateOptimalCount(t *testing.T) {
	generator := NewASCIIArtGenerator()

	tmpDir := t.TempDir()
	testImagePath := filepath.Join(tmpDir, "test.png")

	img := image.NewGray(image.Rect(0, 0, 100, 100))
	f, _ := os.Create(testImagePath)
	png.Encode(f, img)
	f.Close()

	art, _ := generator.GenerateFromImage(testImagePath, 40)

	count := generator.CalculateOptimalCount(art)
	if count != art.LineCount() {
		t.Errorf("CalculateOptimalCount() = %v, want %v", count, art.LineCount())
	}
}
