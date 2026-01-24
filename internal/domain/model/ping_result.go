package model

import (
	"net"
	"time"
)

type PingPacket struct {
	Seq     int
	Nbytes  int
	IPAddr  net.IP
	TTL     int
	Rtt     time.Duration
	ArtLine string
}

func NewPingPacket(seq, nbytes, ttl int, ipAddr net.IP, rtt time.Duration, artLine string) *PingPacket {
	return &PingPacket{
		Seq:     seq,
		Nbytes:  nbytes,
		IPAddr:  ipAddr,
		TTL:     ttl,
		Rtt:     rtt,
		ArtLine: artLine,
	}
}

type PingStatistics struct {
	Addr        string
	PacketsSent int
	PacketsRecv int
	PacketLoss  float64
	MinRtt      time.Duration
	AvgRtt      time.Duration
	MaxRtt      time.Duration
	StdDevRtt   time.Duration
}

func NewPingStatistics(addr string, sent, recv int, loss float64, minRtt, avgRtt, maxRtt, stdDevRtt time.Duration) *PingStatistics {
	return &PingStatistics{
		Addr:        addr,
		PacketsSent: sent,
		PacketsRecv: recv,
		PacketLoss:  loss,
		MinRtt:      minRtt,
		AvgRtt:      avgRtt,
		MaxRtt:      maxRtt,
		StdDevRtt:   stdDevRtt,
	}
}
