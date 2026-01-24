package cli

import (
	"errors"
	"fmt"
	"nyagoPing/internal/application/usecase"
	"nyagoPing/internal/domain/model"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

const (
	ExitCodeOK exitCode = iota
	ExitCodeErrorArgs
	ExitCodeErrorExecution
)

type exitCode int

type Options struct {
	Count          int    `short:"c" long:"count" description:"Pingの送信回数を指定します。"`
	Privilege      bool   `short:"p" long:"privileged" description:"特権モードで実行します。"`
	Version        bool   `short:"v" long:"version" description:"バージョンを表示します。"`
	ASCIIArtPath   string `short:"a" long:"ascii-art" description:"アスキーアートファイルのパスを指定します。" default:".env"`
	Generate       string `short:"g" long:"generate" description:"画像ファイルまたはディレクトリからアスキーアートを生成します。"`
	GenerateOutput string `short:"o" long:"output" description:"生成したアスキーアートの出力先を指定します。" default:".env"`
	GenerateWidth  int    `short:"w" long:"width" description:"生成するアスキーアートの幅を指定します。" default:"80"`
}

type CLI struct {
	pingUseCase     *usecase.PingUseCase
	generateUseCase *usecase.GenerateASCIIArtUseCase
	presenter       *Presenter
	appName         string
	appVersion      string
	appDescription  string
}

func NewCLI(
	pingUseCase *usecase.PingUseCase,
	generateUseCase *usecase.GenerateASCIIArtUseCase,
	presenter *Presenter,
	appName, appVersion, appDescription string,
) *CLI {
	return &CLI{
		pingUseCase:     pingUseCase,
		generateUseCase: generateUseCase,
		presenter:       presenter,
		appName:         appName,
		appVersion:      appVersion,
		appDescription:  appDescription,
	}
}

func (c *CLI) Run(args []string) exitCode {
	code, err := c.run(args)
	if err != nil {
		c.presenter.ShowError(err)
	}
	return code
}

func (c *CLI) run(cliArgs []string) (exitCode, error) {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = c.appName
	parser.Usage = fmt.Sprintf("[オプション...] <ホスト>\n\n%s", c.appDescription)

	args, err := parser.ParseArgs(cliArgs)
	if err != nil {
		if flags.WroteHelp(err) {
			return ExitCodeOK, nil
		}
		return ExitCodeErrorArgs, fmt.Errorf("引数解析エラー: %w", err)
	}

	if opts.Version {
		c.presenter.ShowVersion(c.appName, c.appVersion)
		return ExitCodeOK, nil
	}

	if opts.Generate != "" {
		return c.handleGenerate(&opts)
	}

	if len(args) == 0 {
		return ExitCodeErrorArgs, errors.New("ホスト名を指定してください")
	}
	if len(args) > 1 {
		return ExitCodeErrorArgs, errors.New("ホスト名は1つだけ指定してください")
	}

	return c.handlePing(&opts, args[0])
}

func (c *CLI) handleGenerate(opts *Options) (exitCode, error) {
	outputPath := opts.GenerateOutput
	if outputPath == ".env" {
		execPath, err := os.Executable()
		if err == nil {
			outputPath = filepath.Join(filepath.Dir(execPath), ".env")
		}
	}

	input := &usecase.GenerateInput{
		OutputPath: outputPath,
		Width:      opts.GenerateWidth,
	}

	fileInfo, err := os.Stat(opts.Generate)
	if err != nil {
		return ExitCodeErrorExecution, fmt.Errorf("パスが存在しません: %s", opts.Generate)
	}

	if fileInfo.IsDir() {
		input.ImageDir = opts.Generate
		input.SaveSeparately = false
	} else {
		input.ImagePath = opts.Generate
	}

	output, err := c.generateUseCase.Execute(input)
	if err != nil {
		return ExitCodeErrorExecution, err
	}

	if len(output.Arts) > 0 {
		fmt.Printf("アスキーアートを生成しました (%d個, 各%d行):\n\n", len(output.Arts), output.Arts[0].LineCount())

		if len(output.Arts) > 0 {
			fmt.Printf("=== %s ===\n", output.Filenames[0])
			c.presenter.ShowASCIIArt(output.Arts[0])
		}

		if len(output.Arts) > 1 {
			fmt.Printf("\n他 %d 個の画像も変換されました。\n", len(output.Arts)-1)
		}

		fmt.Printf("\n保存先: %s\n", outputPath)
	}

	return ExitCodeOK, nil
}

func (c *CLI) handlePing(opts *Options, host string) (exitCode, error) {
	count := opts.Count
	autoCount := count == 0

	asciiArtPath := opts.ASCIIArtPath
	if asciiArtPath == ".env" {
		execPath, err := os.Executable()
		if err == nil {
			asciiArtPath = filepath.Join(filepath.Dir(execPath), ".env")
		}
	}

	input := &usecase.PingInput{
		Host:           host,
		Count:          count,
		Privileged:     opts.Privilege,
		ASCIIArtPath:   asciiArtPath,
		AutoCountByArt: autoCount,
	}

	err := c.pingUseCase.Execute(
		input,
		func(packet *model.PingPacket) {
			c.presenter.ShowPingPacket(packet)
		},
		func(stats *model.PingStatistics) {
			c.presenter.ShowPingStatistics(stats)
		},
	)

	if err != nil {
		return ExitCodeErrorExecution, err
	}

	return ExitCodeOK, nil
}

func (c *CLI) Main() {
	code := c.Run(os.Args[1:])
	os.Exit(int(code))
}
