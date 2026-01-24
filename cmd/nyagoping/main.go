package main

import (
	"nyagoPing/internal/application/usecase"
	"nyagoPing/internal/domain/service"
	"nyagoPing/internal/infrastructure/persistence"
	"nyagoPing/internal/infrastructure/ping"
	"nyagoPing/internal/presentation/cli"
)

const (
	appName        = "nyagoping"
	appVersion     = "3.0.0"
	appDescription = "ﾝﾆｬｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞ🐈"
)

func main() {
	pingRepo := ping.NewProBingRepository()
	asciiRepo := persistence.NewFileASCIIArtRepository()
	artGenerator := service.NewASCIIArtGenerator()
	pingUseCase := usecase.NewPingUseCase(pingRepo, asciiRepo, artGenerator)
	generateUseCase := usecase.NewGenerateASCIIArtUseCase(asciiRepo, artGenerator)
	presenter := cli.NewPresenter()
	cliApp := cli.NewCLI(
		pingUseCase,
		generateUseCase,
		presenter,
		appName,
		appVersion,
		appDescription,
	)
	cliApp.Main()
}
