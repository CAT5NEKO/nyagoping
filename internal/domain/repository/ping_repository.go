package repository

import "nyagoPing/internal/domain/model"

type PingRepository interface {
	Ping(target *model.PingTarget, config *model.PingConfig, art *model.ASCIIArt, onRecv func(*model.PingPacket), onFinish func(*model.PingStatistics)) error
}
