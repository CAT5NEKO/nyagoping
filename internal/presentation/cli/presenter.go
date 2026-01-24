package cli

import (
	"fmt"
	"nyagoPing/internal/domain/model"

	"github.com/fatih/color"
)

type Presenter struct{}

func NewPresenter() *Presenter {
	return &Presenter{}
}

func (p *Presenter) ShowPingStart(host string, ip string) {
	fmt.Printf("PING %s (%s)\n", host, ip)
}

func (p *Presenter) ShowPingPacket(packet *model.PingPacket) {
	fmt.Printf("%s %v\n",
		packet.ArtLine,
		color.New(color.FgBlue, color.Bold).Sprint(packet.Rtt),
	)
}

func (p *Presenter) ShowPingStatistics(stats *model.PingStatistics) {
	fmt.Fprintf(color.Output, "\n--- %s 統計 ---\n", stats.Addr)
	fmt.Fprintf(color.Output, "%d送信, %d受信, %.1f%%ロス, avg=%v\n",
		stats.PacketsSent,
		stats.PacketsRecv,
		stats.PacketLoss,
		color.New(color.FgCyan, color.Bold).Sprint(stats.AvgRtt),
	)
}

func (p *Presenter) ShowASCIIArt(art *model.ASCIIArt) {
	for _, line := range art.Lines() {
		fmt.Println(line)
	}
}

func (p *Presenter) ShowError(err error) {
	fmt.Fprintf(color.Output, "[%v] %v\n",
		color.New(color.FgRed, color.Bold).Sprint("ERROR"),
		err,
	)
}

func (p *Presenter) ShowVersion(appName, version string) {
	fmt.Printf("%s version %s\n", appName, version)
}

func (p *Presenter) ShowHelp(usage string) {
	fmt.Println(usage)
}
