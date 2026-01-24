package usecase

import (
	"fmt"
	"nyagoPing/internal/domain/model"
	"nyagoPing/internal/domain/repository"
	"nyagoPing/internal/domain/service"
)

type GenerateASCIIArtUseCase struct {
	asciiRepo    repository.ASCIIArtRepository
	artGenerator *service.ASCIIArtGenerator
}

func NewGenerateASCIIArtUseCase(
	asciiRepo repository.ASCIIArtRepository,
	artGenerator *service.ASCIIArtGenerator,
) *GenerateASCIIArtUseCase {
	return &GenerateASCIIArtUseCase{
		asciiRepo:    asciiRepo,
		artGenerator: artGenerator,
	}
}

type GenerateInput struct {
	ImagePath      string
	ImageDir       string
	OutputPath     string
	Width          int
	SaveSeparately bool
}

type GenerateOutput struct {
	Arts      []*model.ASCIIArt
	Filenames []string
}

func (uc *GenerateASCIIArtUseCase) Execute(input *GenerateInput) (*GenerateOutput, error) {
	var arts []*model.ASCIIArt
	var filenames []string

	if input.ImageDir != "" {
		generatedArts, generatedFilenames, err := uc.artGenerator.GenerateFromImagesInDirectory(input.ImageDir, input.Width)
		if err != nil {
			return nil, fmt.Errorf("ディレクトリからのアスキーアート生成エラー: %w", err)
		}
		arts = generatedArts
		filenames = generatedFilenames

		if !input.SaveSeparately && len(arts) > 0 {
			if err := uc.asciiRepo.Save(input.OutputPath, arts[0]); err != nil {
				return nil, fmt.Errorf("アスキーアート保存エラー: %w", err)
			}
		}
	} else if input.ImagePath != "" {
		art, err := uc.artGenerator.GenerateFromImage(input.ImagePath, input.Width)
		if err != nil {
			return nil, fmt.Errorf("画像からのアスキーアート生成エラー: %w", err)
		}
		arts = append(arts, art)
		filenames = append(filenames, input.ImagePath)

		if err := uc.asciiRepo.Save(input.OutputPath, art); err != nil {
			return nil, fmt.Errorf("アスキーアート保存エラー: %w", err)
		}
	} else {
		return nil, fmt.Errorf("画像パスまたはディレクトリパスを指定してください")
	}

	return &GenerateOutput{
		Arts:      arts,
		Filenames: filenames,
	}, nil
}
