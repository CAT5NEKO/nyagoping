package service

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"nyagoPing/internal/domain/model"
	"os"
	"path/filepath"
	"strings"
)

type ASCIIArtGenerator struct{}

func NewASCIIArtGenerator() *ASCIIArtGenerator {
	return &ASCIIArtGenerator{}
}

var asciiChars = []rune{
	' ', '.', '\'', '`', '^', '"', ',', ':', '~', '-', '_', '+',
	'<', '>', 'i', '!', 'l', 'I', '?', '}', '{', '1', ')', '(',
	'|', '\\', '/', 't', 'f', 'j', 'r', 'x', 'n', 'u', 'v', 'c',
	'z', 'X', 'Y', 'U', 'J', 'C', 'L', 'Q', '0', 'O', 'Z', 'm',
	'w', 'q', 'p', 'd', 'b', 'k', 'h', 'a', 'o', '*', '#', 'M',
	'W', '&', '8', '%', 'B', '@', '$',
}

const (
	maxTerminalWidth  = 200
	maxTerminalHeight = 60
)

func (g *ASCIIArtGenerator) GenerateFromImage(imagePath string, width int) (*model.ASCIIArt, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("画像ファイルを開けません: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("画像をデコードできません: %w", err)
	}

	return g.convertImageToASCII(img, width)
}

func (g *ASCIIArtGenerator) GenerateFromImagesInDirectory(dirPath string, width int) ([]*model.ASCIIArt, []string, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("ディレクトリが存在しません: %s", dirPath)
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, nil, fmt.Errorf("ディレクトリを読み込めません: %w", err)
	}

	var arts []*model.ASCIIArt
	var filenames []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		ext := strings.ToLower(filepath.Ext(filename))

		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue
		}

		imagePath := filepath.Join(dirPath, filename)
		art, err := g.GenerateFromImage(imagePath, width)
		if err != nil {
			return nil, nil, fmt.Errorf("画像 %s の変換エラー: %w", filename, err)
		}

		arts = append(arts, art)
		filenames = append(filenames, filename)
	}

	if len(arts) == 0 {
		return nil, nil, fmt.Errorf("ディレクトリ内に画像ファイル(.jpg, .jpeg, .png)が見つかりません")
	}

	return arts, filenames, nil
}

func (g *ASCIIArtGenerator) convertImageToASCII(img image.Image, width int) (*model.ASCIIArt, error) {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	if width <= 0 {
		width = 80
	}

	aspectRatio := float64(imgHeight) / float64(imgWidth)
	height := int(float64(width) * aspectRatio * 0.55)

	if width > maxTerminalWidth {
		width = maxTerminalWidth
		height = int(float64(width) * aspectRatio * 0.55)
	}
	if height > maxTerminalHeight {
		height = maxTerminalHeight
		width = int(float64(height) / aspectRatio / 0.55)
	}

	if width <= 0 {
		width = 10
	}
	if height <= 0 {
		height = 1
	}

	var lines []string

	for y := 0; y < height; y++ {
		var line strings.Builder
		for x := 0; x < width; x++ {
			imgX := int(float64(x) * float64(imgWidth) / float64(width))
			imgY := int(float64(y) * float64(imgHeight) / float64(height))

			r, g, b, _ := img.At(imgX, imgY).RGBA()

			gray := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 257.0

			charIndex := int(gray / 255.0 * float64(len(asciiChars)-1))
			if charIndex >= len(asciiChars) {
				charIndex = len(asciiChars) - 1
			}

			line.WriteRune(asciiChars[charIndex])
		}
		lines = append(lines, line.String())
	}

	return model.NewASCIIArt(lines)
}

func (g *ASCIIArtGenerator) CalculateOptimalCount(art *model.ASCIIArt) int {
	lineCount := art.LineCount()
	if lineCount == 0 {
		return 10
	}
	return lineCount
}
