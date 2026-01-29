# mmdc へのコントリビューション

mmdc へのコントリビューションに興味を持っていただきありがとうございます！このドキュメントでは、コントリビューションのガイドラインと手順を説明します。

## 目次

- [行動規範](#行動規範)
- [はじめに](#はじめに)
- [開発環境のセットアップ](#開発環境のセットアップ)
- [変更の作成](#変更の作成)
- [変更の提出](#変更の提出)
- [スタイルガイドライン](#スタイルガイドライン)
- [テスト](#テスト)
- [ドキュメント](#ドキュメント)

## 行動規範

このプロジェクトは[Contributor Covenant行動規範](CODE_OF_CONDUCT.ja.md)に従います。参加することで、この規範を守ることが期待されます。

## はじめに

### 前提条件

- Go 1.23.2 以降
- Chrome または Chromium ブラウザ
- Git

### フォークとクローン

1. GitHub でリポジトリをフォーク
2. フォークをローカルにクローン:

   ```bash
   git clone https://github.com/YOUR_USERNAME/mmdc.git
   cd mmdc
   ```

3. upstream リモートを追加:

   ```bash
   git remote add upstream https://github.com/ageha734/mmdc.git
   ```

## 開発環境のセットアップ

### 依存関係のインストール

```bash
go mod download
```

### ビルド

```bash
go build -o mmdc ./cmd
```

### テストの実行

```bash
go test -v ./...
```

### 開発ツールのインストール

```bash
# golangci-lint のインストール
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# mockgen のインストール（モック生成用）
go install go.uber.org/mock/mockgen@latest
```

## 変更の作成

### ブランチ命名規則

変更用のブランチを作成:

```bash
git checkout -b <type>/<description>
```

タイプ:

- `feature/` - 新機能
- `fix/` - バグ修正
- `docs/` - ドキュメント変更
- `refactor/` - コードリファクタリング
- `test/` - テストの追加・更新
- `chore/` - メンテナンスタスク

例:

- `feature/add-webp-support`
- `fix/markdown-extraction-bug`
- `docs/improve-readme`

### コミットメッセージ

[Conventional Commits](https://www.conventionalcommits.org/ja/)に従います:

```text
<type>(<scope>): <description>
```

タイプ:

- `feat`: 新機能
- `fix`: バグ修正
- `chore`: メンテナンス

例:

```text
feat(renderer): WebP出力形式のサポートを追加
fix(markdown): ネストされたコードブロックを正しく処理
```

## 変更の提出

### プルリクエストのプロセス

1. 最新のupstream変更でブランチを更新:

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. すべてのテストとリンターを実行:

   ```bash
   go test -v ./...
   golangci-lint run
   ```

3. 変更をプッシュ:

   ```bash
   git push origin <your-branch>
   ```

4. GitHub でプルリクエストを作成
5. PRテンプレートを完全に記入

### PRの要件

- [ ] すべてのテストが通過
- [ ] コードがスタイルガイドラインに従っている
- [ ] ドキュメントが更新されている（該当する場合）
- [ ] CHANGELOG.md が更新されている
- [ ] コミットメッセージがConventional Commitsに従っている
- [ ] マージコンフリクトがない

## スタイルガイドライン

### Goコードスタイル

- 公式の[Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)に従う
- フォーマットには `gofmt` を使用（golangci-lintで自動実行）
- リンティングには `golangci-lint` を使用:

```bash
golangci-lint run
```

### コード構成

```text
mmdc/
├── cmd/              # CLIエントリーポイント
├── config/           # 設定
├── domain/           # ドメインモデル
├── handler/          # HTTP/IOハンドラー
├── infra/            # インフラストラクチャ（外部サービス）
├── pkg/              # 再利用可能なパッケージ
├── test/mock/        # テストモック
└── usecase/          # ビジネスロジック
```

### 命名規則

- 説明的な名前を使用
- 略語は一貫した大文字小文字で（例: `HTML`, `URL`, `ID`）
- インターフェース名は動作を説明
- 短いループ以外では1文字の変数名を避ける

### エラーハンドリング

- エラーは常に明示的に処理
- ドメインエラーにはカスタムエラー型を使用
- `fmt.Errorf("...: %w", err)` でコンテキストを付けてラップ

```go
// 良い例
if err != nil {
    return fmt.Errorf("ダイアグラムのレンダリングに失敗: %w", err)
}

// 避ける
if err != nil {
    return err
}
```

## テスト

### テスト構造

- テストはテスト対象のコードと同じパッケージに配置
- 複数のテストケースにはテーブル駆動テストを使用
- 意味のあるテスト名を使用

```go
func TestMermaidDiagram_OutputFormat(t *testing.T) {
    tests := []struct {
        name    string
        format  string
        want    OutputFormat
        wantErr bool
    }{
        {"有効なsvg", "svg", OutputFormatSVG, false},
        {"有効なpng", "png", OutputFormatPNG, false},
        {"無効な形式", "invalid", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

### テストの実行方法

```bash
# すべてのテストを実行
go test -v ./...

# カバレッジ付きでテスト実行
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 特定パッケージのテストを実行
go test -v ./pkg/markdown/...

# レース検出付きでテスト実行
go test -race ./...
```

### モック

mockgenを使用してモックを生成:

```bash
go generate ./...
```

## ドキュメント

### コードドキュメント

- エクスポートされるすべての型、関数、メソッドにドキュメントを記載
- コメントには完全な文を使用
- 複雑な機能には例を含める

```go
// RenderDiagram は Mermaid ダイアグラムを指定された出力形式でレンダリングします。
// ヘッドレス Chrome ブラウザを使用してダイアグラムをレンダリングします。
//
// 例:
//
//     output, err := RenderDiagram(ctx, "graph TD; A-->B", FormatSVG)
func RenderDiagram(ctx context.Context, code string, format OutputFormat) ([]byte, error)
```

### README の更新

新機能を追加する際:

1. 機能リストを更新
2. 使用例を追加
3. フラグが変更された場合はコマンドリファレンスを更新

### CHANGELOG の更新

「Unreleased」セクションにエントリを追加:

```markdown
## [Unreleased]

### Added
- 新機能の説明

### Changed
- 変更された動作の説明

### Fixed
- バグ修正の説明
```

## 質問がありますか？

コントリビューションについて質問がある場合:

1. 既存のissueとdiscussionを確認
2. 新しいdiscussionまたはissueを作成
3. メンテナーに連絡

コントリビューションありがとうございます！
