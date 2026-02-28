-- 1. 地理空間計算用の拡張機能を有効化
CREATE EXTENSION IF NOT EXISTS postgis;

-- 2. ユーザーテーブル
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. 店舗テーブル (Spot Entity)
CREATE TABLE spots (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mesh_id VARCHAR(50) NOT NULL, 
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    -- 修正: 誰がそのメッシュの代表として選んだかを保存するカラムを追加
    registered_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 修正: 「1人につき1メッシュ1軒」という蒸留の制約を付与
    -- これによりリポジトリの ON CONFLICT (mesh_id, registered_user_id) が動作します
    UNIQUE (mesh_id, registered_user_id)
);

-- 4. 投稿テーブル (Post Entity)
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    spot_id INTEGER NOT NULL REFERENCES spots(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL, -- 修正: 表示の高速化のため投稿時のユーザー名を非正規化して保持
    image_url TEXT,
    caption TEXT NOT NULL,
    posted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 5. インデックスの作成
CREATE INDEX idx_spots_location ON spots USING GIST (location);
CREATE INDEX idx_spots_mesh_id ON spots (mesh_id);
CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_posts_spot_id ON posts (spot_id);
-- 修正: 共鳴者検索（FindResonantUsers）を高速化するための複合インデックス
CREATE INDEX idx_spots_user_mesh ON spots (registered_user_id, mesh_id);

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