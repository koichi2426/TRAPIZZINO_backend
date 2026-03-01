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
    -- 最初にこのスポット（メッシュ）を登録したユーザーを保持
    registered_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 【修正】一意制約を mesh_id 単体に変更
    -- これにより、誰が投稿しても同じ場所なら conflict が発生し、Go側の ON CONFLICT (mesh_id) が動作します
    CONSTRAINT unique_mesh_per_location UNIQUE (mesh_id)
);

-- 4. 投稿テーブル (Post Entity)
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    spot_id INTEGER NOT NULL REFERENCES spots(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL, 
    image_url TEXT,
    caption TEXT NOT NULL,
    posted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 5. インデックスの作成
CREATE INDEX idx_spots_location ON spots USING GIST (location);
-- UNIQUE制約（unique_mesh_per_location）により mesh_id のインデックスは自動生成されるため、明示的な作成は不要
CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_posts_spot_id ON posts (spot_id);
-- 共鳴者検索（FindResonantUsers）を高速化するための複合インデックス
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