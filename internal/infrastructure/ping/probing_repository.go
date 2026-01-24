package ping

import (
	"fmt"
	"nyagoPing/internal/domain/model"
	"nyagoPing/internal/domain/repository"
	"os"
	"os/signal"
	"runtime"

	probing "github.com/prometheus-community/pro-bing"
)

type ProBingRepository struct{}

func NewProBingRepository() repository.PingRepository {
	return &ProBingRepository{}
}

func (r *ProBingRepository) Ping(
	target *model.PingTarget,
	config *model.PingConfig,
	art *model.ASCIIArt,
	onRecv func(*model.PingPacket),
	onFinish func(*model.PingStatistics),
) error {
	pinger, err := probing.NewPinger(target.Host())
	if err != nil {
		return fmt.Errorf("Pingerの初期化エラー: %w", err)
	}

	pinger.Count = config.Count()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pinger.Stop()
	}()

	if pinger.IPAddr() != nil {
		target.SetIP(pinger.IPAddr().IP)
	}

	pinger.OnRecv = func(pkt *probing.Packet) {
		artLine := art.GetLineBySeq(pkt.Seq)
		packet := model.NewPingPacket(
			pkt.Seq,
			pkt.Nbytes,
			pkt.TTL,
			pkt.IPAddr.IP,
			pkt.Rtt,
			artLine,
		)
		onRecv(packet)
		
		if pkt.Seq >= art.LineCount()-1 {
			pinger.Stop()
		}
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		statistics := model.NewPingStatistics(
			stats.Addr,
			stats.PacketsSent,
			stats.PacketsRecv,
			stats.PacketLoss,
			stats.MinRtt,
			stats.AvgRtt,
			stats.MaxRtt,
			stats.StdDevRtt,
		)
		onFinish(statistics)
	}

	if config.Privileged() || runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	if err := pinger.Run(); err != nil {
		return fmt.Errorf("Ping実行エラー: %w", err)
	}

	return nil
}
