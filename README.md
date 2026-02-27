## 🚀 セットアップ手順

### 0. 必要なツールのインストール

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Atlas CLI](https://atlasgo.io/)

macOSの場合（Homebrew推奨）:
```sh
brew install ariga/tap/atlas
```

### 1. 環境変数の準備

`app/` ディレクトリ内に `.env` ファイルを作成し、データベース接続情報などを記述します。

```bash
# テンプレートをコピーして編集
cp app/.env.example app/.env

```

### 2. コンテナのビルドと起動

Docker Compose を使用して、Go アプリケーションとデータベースのコンテナを立ち上げます。

```bash
# ルートディレクトリで実行
docker-compose up -d --build

```

### 3. データベースマイグレーション

Atlas を使用して、`migrations/` 内の SQL スキーマをデータベースに適用します。

```bash
# スキーマを適用（URLは環境に合わせて調整）
atlas migrate apply \
  --dir "file://migrations" \
  --url "postgres://user:pass@localhost:5432/dbname?sslmode=disable"

```

### 4. 動作確認

サーバーが正常に起動しているか、ヘルスチェックエンドポイントを叩いて確認します。

```bash
curl http://localhost:8000/health
# {"status": "ok"} と返ってくれば成功

```
