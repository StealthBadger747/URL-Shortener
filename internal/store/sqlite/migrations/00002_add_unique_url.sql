-- +goose Up
DELETE FROM urls
WHERE rowid NOT IN (
  SELECT MIN(rowid)
  FROM urls
  GROUP BY url
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_urls_unique_url ON urls(url);

-- +goose Down
DROP INDEX IF EXISTS idx_urls_unique_url;
