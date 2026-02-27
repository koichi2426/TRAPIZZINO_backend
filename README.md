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

Atlas を使用して、`app/migrations/` 内の SQL スキーマをデータベースに適用します。

> ⚠️ **DB接続URLのユーザー名・パスワード・DB名などは、必ず `app/.env` または `app/.env.example` の値を参照して設定してください。**
> 例：
>   DB_USER=user
>   DB_PASSWORD=password
>   DB_NAME=trapizzino
> の場合、
>   --url "postgres://user:password@localhost:5432/trapizzino?sslmode=disable"
> となります。

#### ✅ チェックサムエラーが出た場合

Atlasはマイグレーションファイルの整合性を保つため、チェックサムファイル（atlas.sum）を利用します。
もし `checksum file not found` や `checksum error` が出た場合は、下記コマンドでチェックサムファイルを再生成してください。

```bash
atlas migrate hash --dir "file://app/migrations"
```

このコマンドは app/migrations ディレクトリ内のマイグレーションファイルのハッシュ情報を生成し、atlas.sumファイルを作成します。

```bash
# スキーマを適用（--urlは.envの値に合わせて調整）
atlas migrate apply \
  --dir "file://app/migrations" \
  --url "postgres://user:password@localhost:5432/trapizzino?sslmode=disable" \
  --allow-dirty
```

### 4. 動作確認

サーバーが正常に起動しているか、ヘルスチェックエンドポイントを叩いて確認します。

```bash
curl http://localhost:8000/health
# {"status": "ok"} と返ってくれば成功

```
