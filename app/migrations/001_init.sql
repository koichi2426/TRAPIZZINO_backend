-- ==============================================================================
-- 001_init.sql: 初期スキーマ定義
-- PostGISの有効化と店舗・投稿・ユーザー管理テーブルの構築
-- ==============================================================================

-- 1. 地理空間計算用の拡張機能を有効化
CREATE EXTENSION IF NOT EXISTS postgis;

-- 2. ユーザーテーブル (User Entity)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. 店舗テーブル (Spot Entity)
-- 「情報の蒸留」をDBレベルで強制するため、mesh_idにUNIQUE制約を付与
CREATE TABLE spots (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mesh_id VARCHAR(50) NOT NULL UNIQUE, -- 1つのメッシュには「最適解」として1軒のみ保持
    location GEOGRAPHY(POINT, 4326) NOT NULL, -- 緯度経度データ（SRID 4326）
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 4. 投稿テーブル (Post Entity)
-- SpotとUserを紐付ける個人の記録
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    spot_id INTEGER NOT NULL REFERENCES spots(id) ON DELETE CASCADE,
    image_url TEXT, -- Firebase Storage等のURL（null許容）
    caption TEXT NOT NULL,
    posted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 5. インデックスの作成
-- 空間検索（周辺検索や激戦区計算）を高速化するためのGISTインデックス
CREATE INDEX idx_spots_location ON spots USING GIST (location);

-- 特定のメッシュIDでの検索を高速化
CREATE INDEX idx_spots_mesh_id ON spots (mesh_id);

-- ユーザーごとの投稿一覧取得を高速化
CREATE INDEX idx_posts_user_id ON posts (user_id);

-- 特定スポットに紐づく投稿の取得を高速化
CREATE INDEX idx_posts_spot_id ON posts (spot_id);

-- 6. 更新日時自動更新用の関数とトリガー
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_spots_updated_at BEFORE UPDATE ON spots FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();