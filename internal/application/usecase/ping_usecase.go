package usecase

import (
	"fmt"
	"nyagoPing/internal/domain/model"
	"nyagoPing/internal/domain/repository"
	"nyagoPing/internal/domain/service"
)

type PingUseCase struct {
	pingRepo    repository.PingRepository
	asciiRepo   repository.ASCIIArtRepository
	artGenerator *service.ASCIIArtGenerator
}

func NewPingUseCase(
	pingRepo repository.PingRepository,
	asciiRepo repository.ASCIIArtRepository,
	artGenerator *service.ASCIIArtGenerator,
) *PingUseCase {
	return &PingUseCase{
		pingRepo:     pingRepo,
		asciiRepo:    asciiRepo,
		artGenerator: artGenerator,
	}
}

type PingInput struct {
	Host           string
	Count          int
	Privileged     bool
	ASCIIArtPath   string
	AutoCountByArt bool
}

func (uc *PingUseCase) Execute(
	input *PingInput,
	onRecv func(*model.PingPacket),
	onFinish func(*model.PingStatistics),
) error {
	target, err := model.NewPingTarget(input.Host)
	if err != nil {
		return fmt.Errorf("ターゲット作成エラー: %w", err)
	}

	art, err := uc.asciiRepo.Load(input.ASCIIArtPath)
	if err != nil {
		return fmt.Errorf("アスキーアート読み込みエラー: %w", err)
	}

	count := input.Count
	if input.AutoCountByArt {
		count = uc.artGenerator.CalculateOptimalCount(art)
	}

	config, err := model.NewPingConfig(count, input.Privileged)
	if err != nil {
		return fmt.Errorf("設定作成エラー: %w", err)
	}

	return uc.pingRepo.Ping(target, config, art, onRecv, onFinish)
}
