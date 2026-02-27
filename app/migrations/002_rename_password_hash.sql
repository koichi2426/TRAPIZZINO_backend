-- 002_rename_password_hash.sql: usersテーブルのカラム名変更
ALTER TABLE users RENAME COLUMN password_hash TO hashed_password;
