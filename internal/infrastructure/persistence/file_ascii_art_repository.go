package persistence

import (
	"bufio"
	"fmt"
	"nyagoPing/internal/domain/model"
	"nyagoPing/internal/domain/repository"
	"os"
)

type FileASCIIArtRepository struct{}

func NewFileASCIIArtRepository() repository.ASCIIArtRepository {
	return &FileASCIIArtRepository{}
}

func (r *FileASCIIArtRepository) Load(path string) (*model.ASCIIArt, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ファイルを開けません: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	return model.NewASCIIArt(lines)
}

func (r *FileASCIIArtRepository) Save(path string, art *model.ASCIIArt) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ファイルを作成できません: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range art.Lines() {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("ファイル書き込みエラー: %w", err)
		}
	}

	return writer.Flush()
}
