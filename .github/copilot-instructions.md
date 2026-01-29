# Copilot Code Review Instructions

このドキュメントはGitHub Copilotがコードレビューを行う際の指針を定義します。

## プロジェクト概要

- **プロジェクト名**: mmdc (Mermaid Diagram CLI)
- **言語**: Go 1.21+
- **アーキテクチャ**: クリーンアーキテクチャ
- **主要な依存関係**: chromedp (ブラウザ制御)

## アーキテクチャルール

```text
cmd/        → config, domain, handler, infra, usecase, pkg
usecase/    → domain, config, infra, pkg
handler/    → domain, config, pkg
infra/      → domain, config, handler, pkg
domain/     → pkg
config/     → pkg
pkg/        → (依存なし)
```

### 依存関係の違反チェック

以下のインポートパターンは**エラー**として報告してください：

- `domain` が `usecase`, `handler`, `infra` をインポート
- `pkg` が他のパッケージをインポート
- `config` が `domain` 以外の内部パッケージをインポート
- 循環依存

## コードスタイル

### 必須

- [ ] `gofmt -s` でフォーマット済み
- [ ] `goimports` でインポート整理済み
- [ ] エクスポートされた関数/型にはGoDocコメント
- [ ] エラーは適切にラップ (`fmt.Errorf("context: %w", err)`)
- [ ] コンテキストは最初の引数 (`func Foo(ctx context.Context, ...)`)

### 禁止

- [ ] `panic()` の使用（テスト以外）
- [ ] グローバル変数の変更
- [ ] `init()` 関数での副作用
- [ ] ハードコードされたファイルパス
- [ ] 未使用のコード

## セキュリティチェック

### 重大 (Critical)

- [ ] コマンドインジェクション (`os/exec` の引数)
- [ ] パストラバーサル (ユーザー入力のファイルパス)
- [ ] 機密情報のハードコード

### 警告 (Warning)

- [ ] エラーメッセージでの機密情報漏洩
- [ ] 不適切なファイルパーミッション
- [ ] タイムアウトのないHTTPクライアント

## テストガイドライン

### テスト必須

- [ ] テーブル駆動テスト
- [ ] エラーケースのテスト
- [ ] テストヘルパーは `t.Helper()` を呼び出す
- [ ] 並列実行可能なテストは `t.Parallel()`

### テストファイルの配置

```text
foo.go          → foo_test.go (同一パッケージ)
test/e2e/       → E2Eテスト
test/integration/ → 統合テスト (//go:build integration)
```

## レビュー観点

### パフォーマンス

- [ ] 不要なメモリアロケーション
- [ ] ループ内でのスライス拡張
- [ ] 大きな構造体の値渡し（ポインタを使用すべき）
- [ ] defer のループ内使用

### 可読性

- [ ] 関数は50行以下が望ましい
- [ ] ネストは3レベル以下
- [ ] 変数名は意図を表現
- [ ] マジックナンバーは定数化

### エラーハンドリング

- [ ] エラーは無視しない (`_ = err` は禁止)
- [ ] エラーメッセージは小文字で開始
- [ ] センチネルエラーは `errors.Is()` で比較
- [ ] カスタムエラーは `Error()` を実装

## コメントフォーマット

レビューコメントは以下の形式で記述してください：

```text
[severity]: message

説明（必要に応じて）

推奨修正:
（コード例は別のコードブロックで記述）
```

コード例：

```go
// 修正コード
```

### Severity レベル

| レベル          | 説明                                           |
| --------------- | ---------------------------------------------- |
| `[critical]`    | セキュリティ問題、データ損失の可能性           |
| `[error]`       | バグ、アーキテクチャ違反                       |
| `[warning]`     | パフォーマンス問題、ベストプラクティス違反     |
| `[suggestion]`  | 改善提案、コードスタイル                       |
| `[question]`    | 意図の確認、ディスカッション                   |

## 無視するファイル

以下のファイルはレビュー対象外です：

- `vendor/`
- `*.pb.go` (Protocol Buffers生成ファイル)
- `*_gen.go` (コード生成ファイル)
- `test/e2e/mock/` (テストモックファイル)

## 自動解決条件

以下の条件で指摘を自動的にresolveしてください：

1. 指摘された行が削除された
2. 指摘された問題が修正された
3. 関連するコードが大幅に変更された

## 追加コンテキスト

### プロジェクト固有のパターン

```go
// エラーラッピングのパターン
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// オプションパターン
type Option func(*Config)

func WithTimeout(d time.Duration) Option {
    return func(c *Config) {
        c.Timeout = d
    }
}

// インターフェースの定義場所
// 使用側（usecase）で定義、実装側（infra）で実装
```

### Chromedpに関する注意

- コンテキストのタイムアウトを必ず設定
- ブラウザリソースは適切にクリーンアップ
- サンドボックスモードの考慮（CI環境）
