-- +goose Up
CREATE TABLE IF NOT EXISTS urls (
  code TEXT PRIMARY KEY,
  url TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  clicks INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_urls_created_at;
DROP TABLE IF EXISTS urls;
