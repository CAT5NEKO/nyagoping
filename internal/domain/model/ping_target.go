package model

import (
	"fmt"
	"net"
)

type PingTarget struct {
	host string
	ip   net.IP
}

func NewPingTarget(host string) (*PingTarget, error) {
	if host == "" {
		return nil, fmt.Errorf("ホスト名が空です")
	}
	return &PingTarget{
		host: host,
	}, nil
}

func (pt *PingTarget) Host() string {
	return pt.host
}

func (pt *PingTarget) IP() net.IP {
	return pt.ip
}

func (pt *PingTarget) SetIP(ip net.IP) {
	pt.ip = ip
}
