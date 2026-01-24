//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	appName   = "nyagoping"
	binDir    = "bin"
	mainPath  = "./cmd/nyagoping"
	coverDir  = "coverage"
	coverFile = "coverage/coverage.out"
	coverHTML = "coverage/coverage.html"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "build":
		build()
	case "cross-build":
		crossBuild()
	case "build-windows":
		buildWindows()
	case "build-linux":
		buildLinux()
	case "build-mac":
		buildMac()
	case "clean":
		clean()
	case "test":
		test()
	case "test-unit":
		testUnit()
	case "test-integration":
		testIntegration()
	case "coverage":
		coverage()
	case "fmt":
		format()
	case "deps":
		deps()
	case "run":
		run()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("nyagoping ビルドスクリプト")
	fmt.Println()
	fmt.Println("使い方:")
	fmt.Println("  go run build.go <command>")
	fmt.Println()
	fmt.Println("コマンド:")
	fmt.Println("  build           ビルドを実行")
	fmt.Println("  cross-build     全プラットフォーム向けにクロスビルド")
	fmt.Println("  build-windows   Windows向けビルド (64bit + 32bit)")
	fmt.Println("  build-linux     Linux向けビルド (64bit + 32bit)")
	fmt.Println("  build-mac       macOS向けビルド (Intel + Apple Silicon)")
	fmt.Println("  clean           ビルドしたやつを削除")
	fmt.Println("  test            全てのテストを実行")
	fmt.Println("  test-unit       ユニットテストのみ実行")
	fmt.Println("  test-integration 統合テストのみ実行")
	fmt.Println("  coverage        カバレッジを取得")
	fmt.Println("  fmt             コードフォーマット")
	fmt.Println("  deps            依存関係をインストール")
	fmt.Println("  run             アプリケーションを実行")
	fmt.Println("  help            ヘルプ")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  go run build.go build")
	fmt.Println("  go run build.go test")
}

func build() {
	fmt.Println("ビルド中...")

	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Printf("エラー: binディレクトリの作成に失敗: %v\n", err)
		os.Exit(1)
	}

	binaryName := appName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	outputPath := filepath.Join(binDir, binaryName)

	cmd := exec.Command("go", "build", "-o", outputPath, mainPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("エラー: ビルドに失敗: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ビルド完了: %s\n", outputPath)
}

func crossBuild() {
	fmt.Println("全プラットフォーム向けクロスビルド中...")

	platforms := []struct {
		goos   string
		goarch string
	}{
		{"windows", "amd64"},
		{"windows", "386"},
		{"linux", "amd64"},
		{"linux", "386"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
	}

	for _, platform := range platforms {
		buildForPlatform(platform.goos, platform.goarch)
	}

	fmt.Println("クロスビルド完了!")
}

func buildForPlatform(goos, goarch string) {
	fmt.Printf("%s/%s 向けビルド中...\n", goos, goarch)

	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Printf("エラー: binディレクトリの作成に失敗: %v\n", err)
		os.Exit(1)
	}

	binaryName := fmt.Sprintf("%s-%s-%s", appName, goos, goarch)
	if goos == "windows" {
		binaryName += ".exe"
	}

	outputPath := filepath.Join(binDir, binaryName)

	cmd := exec.Command("go", "build", "-o", outputPath, mainPath)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("エラー: %s/%s ビルドに失敗: %v\n", goos, goarch, err)
		os.Exit(1)
	}

	fmt.Printf("%s/%s ビルド完了: %s\n", goos, goarch, outputPath)
}

func buildWindows() {
	fmt.Println("Windows向けビルド中...")
	buildForPlatform("windows", "amd64")
	buildForPlatform("windows", "386")
}

func buildLinux() {
	fmt.Println("Linux向けビルド中...")
	buildForPlatform("linux", "amd64")
	buildForPlatform("linux", "386")
}

func buildMac() {
	fmt.Println("macOS向けビルド中...")
	buildForPlatform("darwin", "amd64")
	buildForPlatform("darwin", "arm64")
}

func clean() {
	fmt.Println("クリーニング中...")

	if err := os.RemoveAll(binDir); err != nil {
		fmt.Printf("警告: binディレクトリの削除に失敗: %v\n", err)
	}

	if err := os.RemoveAll(coverDir); err != nil {
		fmt.Printf("警告: coverageディレクトリの削除に失敗: %v\n", err)
	}

	cmd := exec.Command("go", "clean")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	fmt.Println("クリーンアップ完了")
}

func test() {
	fmt.Println("全テストを実行中...")
	runCommand("go", "test", "-v", "./...")
}

func testUnit() {
	fmt.Println("ユニットテストを実行中...")
	runCommand("go", "test", "-v", "./internal/...")
}

func testIntegration() {
	fmt.Println("統合テストを実行中...")
	runCommand("go", "test", "-v", "./test/integration/...")
}

func coverage() {
	fmt.Println("テストカバレッジを取得中...")

	if err := os.MkdirAll(coverDir, 0755); err != nil {
		fmt.Printf("エラー: coverageディレクトリの作成に失敗: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "test", "-coverprofile="+coverFile, "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("エラー: テストに失敗: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("go", "tool", "cover", "-html="+coverFile, "-o", coverHTML)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("エラー: HTMLレポートの生成に失敗: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("カバレッジレポート: %s\n", coverHTML)
}

func format() {
	fmt.Println("コードをフォーマット中...")
	runCommand("go", "fmt", "./...")
}

func deps() {
	fmt.Println("依存関係をインストール中...")
	runCommand("go", "mod", "download")
	runCommand("go", "mod", "tidy")
	fmt.Println("依存関係のインストール完了")
}

func run() {
	fmt.Println("アプリケーションを実行中...")
	args := []string{"run", mainPath + "/main.go"}
	if len(os.Args) > 2 {
		args = append(args, os.Args[2:]...)
	}
	runCommand("go", args...)
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Printf("エラー: コマンドの実行に失敗: %v\n", err)
		os.Exit(1)
	}
}
