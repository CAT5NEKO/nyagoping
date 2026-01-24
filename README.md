# nyagoping

nyagoのpingツール  
Ping時に画像を表示する機能はPinguと変わりませんが、  
画像からAAに変換してそれをPing時に表示する機能など、
色々ほしいなぁと思ってた機能を詰め込んで一から書いてみました。

## インストール


```bash
git clone https://github.com/CAT5NEKO/nyagoping.git
cd nyagoping
go mod tidy
go build -o bin/nyagoping ./cmd/nyagoping

```

# Makefileを使っていい感じにビルド

```bash
make build #単体

make cross-build # 全プラットフォーム

make build-windows   # Windows 64bit + 32bit
make build-linux     # Linux 64bit + 32bit
make build-mac       # macOS Intel + Apple Silicon
```


### Makefileなんてねぇよという方はbuild.goでも似たようなこと実行できます

```bash
go run main.go help # コマンドでやれることはここに記載しています
go run build.go cross-build
go run build.go build-windows
```

## くいっくすたぁと

```bash
nyagoping example.tld
nyagoping -c 5 example.tld               # 固定5回
nyagoping -a myart.txt example.tld       # カスタムAAを使用
nyagoping -g image.png -o myart.txt     # 画像からAA生成
```

## コマンドラインオプション

### PINGする場合

| オプション | 短縮 | 説明 | デフォルト |
|-----------|------|------|-----------|
| --count | -c | Ping送信回数 | 10 |
| --privileged | -p | 特権モード | false |
| --version | -v | バージョン表示 | - |

### カスタムAA使ってPINGする場合
カスタム画像の調整などはAAディレクトリ内部をご確認ください。

| オプション | 短縮 | 説明 | デフォルト |
|-----------|------|------|-----------|
| --ascii-art | -a | AAファイルパス | .env |
| --generate | -g | 画像からAA生成 | - |
| --output | -o | AA出力先 | .env |
| --width | -w | AA幅 | 80 |


