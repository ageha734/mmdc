# mmdc

[![CI](https://github.com/ageha734/mmdc/actions/workflows/ci.yml/badge.svg)](https://github.com/ageha734/mmdc/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ageha734/mmdc)](https://goreportcard.com/report/github.com/ageha734/mmdc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/ageha734/mmdc.svg)](https://pkg.go.dev/github.com/ageha734/mmdc)

[English](README.md)

MermaidダイアグラムをSVG、PNG、PDF形式にレンダリングするGoベースのCLIツールです。chromedpを使用したヘッドレスブラウザレンダリングにより、公式の[@mermaid-js/mermaid-cli](https://github.com/mermaid-js/mermaid-cli)のネイティブGo代替を提供します。

## 特徴

- 🎨 **複数の出力形式**: SVG、PNG、PDF対応
- 📝 **Markdownサポート**: Markdownファイルに埋め込まれたMermaidダイアグラムを抽出・レンダリング
- 🎭 **テーマサポート**: default、forest、dark、neutralテーマ
- ⚡ **高速レンダリング**: chromedpによる効率的なヘッドレスブラウザレンダリング
- 📦 **シングルバイナリ**: Node.jsやnpmは不要
- 🔧 **高いカスタマイズ性**: サイズ、色、スケールなど細かく設定可能

## インストール

### Go を使用

```bash
go install github.com/ageha734/mmdc/cmd@latest
```

### Homebrew を使用（macOS/Linux）

```bash
brew tap ageha734/mmdc
brew install mmdc
```

### curl を使用（Linux/macOS）

```bash
curl -fsSL https://raw.githubusercontent.com/ageha734/mmdc/master/install.sh | bash
```

### PowerShell を使用（Windows）

```powershell
irm https://raw.githubusercontent.com/ageha734/mmdc/master/install.ps1 | iex
```

### バイナリダウンロード

[Releases](https://github.com/ageha734/mmdc/releases)ページから最新のバイナリをダウンロードできます。

## クイックスタート

### MermaidファイルをSVGにレンダリング

```bash
mmdc -i diagram.mmd -o diagram.svg
```

### カスタムサイズでPNGにレンダリング

```bash
mmdc -i diagram.mmd -o diagram.png -w 1200 -H 800
```

### 埋め込みダイアグラムを含むMarkdownをレンダリング

```bash
mmdc -i document.md -a ./images/
```

### 標準入力から読み込み

```bash
echo "graph TD; A-->B" | mmdc -i - -o output.svg
```

## コマンドリファレンス

```text
mmdc - Mermaid CLI for Go

使用方法:
  mmdc [フラグ]
  mmdc [コマンド]

フラグ:
  -i, --input string              入力Mermaidファイルまたはmarkdownファイル
  -o, --output string             出力ファイルパス（デフォルト: 入力ファイル名 + .svg）
  -e, --outputFormat string       出力形式: svg、png、またはpdf
  -t, --theme string              テーマ: default、forest、dark、またはneutral（デフォルト "default"）
  -w, --width int                 ページ幅（デフォルト 800）
  -H, --height int                ページ高さ（デフォルト 600）
  -b, --backgroundColor string    背景色（デフォルト "white"）
  -s, --scale int                 スケール係数（デフォルト 1）
  -q, --quiet                     ログ出力を抑制
  -c, --configFile string         Mermaid設定JSONファイル
  -C, --cssFile string            カスタムスタイリング用CSSファイル
  -I, --svgId string              SVG要素ID
  -f, --pdfFit                    PDFをチャートに合わせてスケール
  -a, --artefacts string          出力アーティファクトパス（Markdown入力用）
  -p, --puppeteerConfigFile       Puppeteer設定JSONファイル
      --iconPacks strings         使用するアイコンパック
  -v, --version                   バージョン情報

コマンド:
  render      Mermaidダイアグラムをレンダリング
  help        任意のコマンドのヘルプ
```

## 使用例

### 基本的な使用方法

```bash
# MermaidをSVGに変換
mmdc -i flowchart.mmd

# ダークテーマでPNGに変換
mmdc -i sequence.mmd -o sequence.png -t dark

# PDFに変換
mmdc -i diagram.mmd -o diagram.pdf -e pdf
```

### Markdown処理

```bash
# Markdownからダイアグラムを抽出してimagesディレクトリに保存
mmdc -i README.md -a ./images/

# カスタム出力形式でMarkdownを処理
mmdc -i docs.md -a ./diagrams/ -e png
```

### 高度な設定

```bash
# カスタムMermaid設定を使用
mmdc -i diagram.mmd -c mermaid-config.json

# カスタムCSSスタイリングを適用
mmdc -i diagram.mmd -C custom.css -o styled.svg

# 高解像度PNG出力
mmdc -i diagram.mmd -o hires.png -s 2 -w 1600 -H 1200
```

### パイプライン使用

```bash
# スクリプト出力からダイアグラムを生成
cat <<EOF | mmdc -i - -o pipeline.svg
graph LR
    A[開始] --> B{判定}
    B -->|はい| C[処理]
    B -->|いいえ| D[終了]
EOF
```

## 設定ファイル

Mermaid設定用のJSON設定ファイルを使用できます:

```json
{
  "theme": "forest",
  "themeVariables": {
    "primaryColor": "#ff6b6b",
    "secondaryColor": "#4ecdc4"
  },
  "flowchart": {
    "curve": "basis"
  },
  "sequence": {
    "mirrorActors": false
  }
}
```

使用方法:

```bash
mmdc -i diagram.mmd -c config.json
```

## 要件

- Go 1.23.2以降（ソースからビルドする場合）
- Chrome/Chromiumブラウザ（自動検出）

## アーキテクチャ

```text
mmdc/
├── cmd/              # CLIエントリーポイント（Cobra）
├── config/           # 設定定数
├── domain/           # ドメインモデル（MermaidDiagram、OutputFormat）
├── handler/          # ブラウザとファイル処理
├── infra/            # インフラストラクチャ（renderer、logger、file I/O）
├── pkg/              # ユーティリティパッケージ（markdown、validation、path）
├── test/mock/        # テストモック
└── usecase/          # ビジネスロジック
```

## コントリビューション

コントリビューションを歓迎します！プルリクエストを送信する前に[コントリビューションガイドライン](CONTRIBUTING.ja.md)をお読みください。

1. リポジトリをフォーク
2. フィーチャーブランチを作成（`git checkout -b feature/amazing-feature`）
3. 変更をコミット（`git commit -m 'Add amazing feature'`）
4. ブランチにプッシュ（`git push origin feature/amazing-feature`）
5. プルリクエストを作成

## セキュリティ

セキュリティに関する問題については、[セキュリティポリシー](SECURITY.ja.md)をご覧ください。

## ライセンス

このプロジェクトはMITライセンスの下でライセンスされています。詳細は[LICENSE](LICENSE.md)ファイルを参照してください。

## 謝辞

- [mermaid-js/mermaid](https://github.com/mermaid-js/mermaid) - 素晴らしいダイアグラムライブラリ
- [chromedp/chromedp](https://github.com/chromedp/chromedp) - Go用Chrome DevTools Protocol
- [spf13/cobra](https://github.com/spf13/cobra) - CLIフレームワーク

## 関連プロジェクト

- [@mermaid-js/mermaid-cli](https://github.com/mermaid-js/mermaid-cli) - 公式Node.jsベースのMermaid CLI
