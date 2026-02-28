## 🚀 セットアップ手順

### 0. 必要なツールのインストール

以下のツールがインストールされているか確認してください。

* **Docker / Docker Compose**
* **Atlas CLI**

```bash
# Linux (VPS) の場合
curl -sSf https://atlasgo.sh | sh

```

### 1. ソースコードの取得

VPS上でリポジトリをクローンし、プロジェクトディレクトリへ移動します。

```bash
git clone https://github.com/koichi2426/TRAPIZZINO_backend.git
cd TRAPIZZINO_backend

```

※ すでにクローン済みの場合は最新化してください。

```bash
git pull origin main

```

### 2. 環境変数の準備

`app/` ディレクトリ内に `.env` ファイルを作成します。

```bash
cp app/.env.example app/.env
# Vimで編集
vi app/.env

```

### 3. コンテナのビルドと起動

環境によってコマンドが異なるため、動く方を実行してください。

```bash
# パターンA (最新のDocker)
docker compose up -d --build

# パターンB (旧バージョン)
docker-compose up -d --build

```

### 4. データベースマイグレーション

> ⚠️ **DB接続URLは `app/.env` の値に合わせて設定してください。**

```bash
# チェックサムの整合性確保
atlas migrate hash --dir "file://app/migrations"

# スキーマを適用
atlas migrate apply \
  --dir "file://app/migrations" \
  --url "postgres://user:password@localhost:5432/trapizzino?sslmode=disable" \
  --allow-dirty

```

### 5. 動作確認

正常にデプロイされたか、2つの環境から確認します。

#### A. 内部確認（VPS内から実行）

```bash
curl http://127.0.0.1:8000/health

```

#### B. 外部確認（手元のMacなどから実行）

```bash
# api.example.com はご自身のドメインに読み替えてください
curl https://api.example.com/health

```
