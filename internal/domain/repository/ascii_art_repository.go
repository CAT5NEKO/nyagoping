package repository

import "nyagoPing/internal/domain/model"

type ASCIIArtRepository interface {
	Load(path string) (*model.ASCIIArt, error)
	Save(path string, art *model.ASCIIArt) error
}
