package model

import "fmt"

type PingConfig struct {
	count      int
	privileged bool
}

func NewPingConfig(count int, privileged bool) (*PingConfig, error) {
	if count < 0 {
		return nil, fmt.Errorf("count は0以上である必要があります: %d", count)
	}
	return &PingConfig{
		count:      count,
		privileged: privileged,
	}, nil
}

func (pc *PingConfig) Count() int {
	return pc.count
}

func (pc *PingConfig) Privileged() bool {
	return pc.privileged
}

func (pc *PingConfig) SetCount(count int) error {
	if count < 0 {
		return fmt.Errorf("count は0以上である必要があります: %d", count)
	}
	pc.count = count
	return nil
}
